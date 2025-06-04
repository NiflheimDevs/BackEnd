package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/infrastructure/websocket"
)

type ChatHandler struct {
	validator   *validator.Validate
	JWTService  services.JWT
	Constants   *bootstrap.Constants
	Hub         *websocket.Hub
	ChatService services.ChatService
}

func NewChatHandler(
	validator *validator.Validate,
	jwtService services.JWT,
	constants *bootstrap.Constants,
	hub *websocket.Hub,
	chatService services.ChatService,
) *ChatHandler {
	return &ChatHandler{
		validator:   validator,
		JWTService:  jwtService,
		Constants:   constants,
		Hub:         hub,
		ChatService: chatService,
	}
}

func (ch *ChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomIDStr := chi.URLParam(r, "room_id")
	token := chi.URLParam(r, "token")

	if roomIDStr == "" || token == "" {
		http.Error(w, "Missing room_id or token", http.StatusBadRequest)
		return
	}

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room_id format", http.StatusBadRequest)
		return
	}

	claims, err := ch.JWTService.VerifyToken(token)
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

	conn := r.Context().Value(ch.Constants.Context.WebSocketConnection)

	client := websocket.NewClient(ch.Hub, conn, roomID, userID, &ch.Constants.WebsocketSetting, ch.ChatService)
	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}

func (ch *ChatHandler) GetAllRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ch.Constants.Context.UserID).(int)

	rooms := ch.ChatService.GetAllRoom(userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (ch *ChatHandler) GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	roomIDString := chi.URLParam(r, "room_id")
	roomID, _ := strconv.Atoi(roomIDString)

	userID := r.Context().Value(ch.Constants.Context.UserID).(int)

	messages := ch.ChatService.GetRoomMessages(userID, roomID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (ch *ChatHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	type roomParam struct {
		TargetUserID int `json:"target_user_id" validate:"required"`
	}

	params := Validated[roomParam](ch.validator, r)

	userID := r.Context().Value(ch.Constants.Context.UserID).(int)

	roomInfo := ch.ChatService.CreateUserRoom(userID, params.TargetUserID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(roomInfo); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}
