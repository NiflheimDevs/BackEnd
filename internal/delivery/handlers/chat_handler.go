package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/infrastructure/websocket"
)

type ChatHandler struct {
	validator  *validator.Validate
	JWTService services.JWT
	Constants  *bootstrap.Constants
	Hub        *websocket.Hub
}

func NewChatHandler(
	validator *validator.Validate,
	jwtService services.JWT,
	constants *bootstrap.Constants,
	hub *websocket.Hub,
) *ChatHandler {
	return &ChatHandler{
		validator:  validator,
		JWTService: jwtService,
		Constants:  constants,
		Hub:        hub,
	}
}

func (ch *ChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	type roomConnectionParam struct {
		RoomID uint   `uri:"room_id" validate:"required"`
		Token  string `uri:"token" validate:"required"`
	}

	param := Validated[roomConnectionParam](ch.validator, r)

	claims, err := ch.JWTService.VerifyToken(param.Token)
	if err != nil {
		panic(err)
	}

	if claims == nil {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_TOKEN_EXPIRED,
			},
		})
	}

	userID := uint(claims["sub"].(float64))

	conn := r.Context().Value(ch.Constants.Context.WebSocketConnection)

	client := websocket.NewClient(ch.Hub, conn, param.RoomID, userID, &ch.Constants.WebsocketSetting)
	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
