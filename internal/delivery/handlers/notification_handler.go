package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
	websocket "github.com/niflheimdevs/backend/internal/infrastructure/websocket"
)

type NotificationHandler struct {
	Validator    *validator.Validate
	Constants    *bootstrap.Constants
	JWTService   services.JWT
	Hub          *websocket.Hub
	NotifService services.NotifService
}

func NewNotificationHandler(
	validator *validator.Validate,
	jwtService services.JWT,
	constants *bootstrap.Constants,
	hub *websocket.Hub,
	notifService services.NotifService,
) *NotificationHandler {
	return &NotificationHandler{
		Validator:    validator,
		JWTService:   jwtService,
		Constants:    constants,
		Hub:          hub,
		NotifService: notifService,
	}
}

func (nh *NotificationHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if token == "" {
		http.Error(w, "Missing room_id or token", http.StatusBadRequest)
		return
	}

	claims, err := nh.JWTService.VerifyToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		http.Error(w, "Invalid token payload", http.StatusUnauthorized)
		return
	}
	userID := int(sub)

	conn := r.Context().Value(nh.Constants.Context.WebSocketConnection)

	client := websocket.NewClient(nh.Hub, conn, 0, userID, &nh.Constants.WebsocketSetting, nil, nh.NotifService)
	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
