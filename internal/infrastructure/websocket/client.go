package websocketimpl

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
)

type Client struct {
	WebsocketSetting *bootstrap.WebsocketSetting
	Hub              *Hub
	Conn             *websocket.Conn
	Send             chan []byte
	RoomID           int
	UserID           int
	Mu               sync.Mutex
	Done             chan struct{}
	CloseOnce        sync.Once
	ChatService      services.ChatService
	NotifService     services.NotifService
}

func NewClient(
	hub *Hub, conn any, roomID, userID int,
	websocketSetting *bootstrap.WebsocketSetting,
	chatService services.ChatService,
	notifService services.NotifService,
) *Client {
	wsConn, _ := conn.(*websocket.Conn)
	return &Client{
		WebsocketSetting: websocketSetting,
		Hub:              hub,
		Conn:             wsConn,
		Send:             make(chan []byte, websocketSetting.MessageBufferSize),
		RoomID:           roomID,
		UserID:           userID,
		Done:             make(chan struct{}),
		ChatService:      chatService,
		NotifService:     notifService,
	}
}

func (client *Client) ReadPump() error {
	defer client.CloseConnection()

	client.Conn.SetReadLimit(int64(client.WebsocketSetting.MaxMessageSize))
	client.Conn.SetReadDeadline(time.Now().Add(client.WebsocketSetting.ReadTimeout))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(client.WebsocketSetting.ReadTimeout))
		return nil
	})

	for {
		_, rawmessage, err := client.Conn.ReadMessage()
		if err != nil {
			return err
		}

		var message Message

		if err := json.Unmarshal(rawmessage, &message); err != nil {
			continue
		}
		message.Client = client
		message.Timestamp = time.Now()
		message.RoomID = client.RoomID
		switch message.Type {
		case MessageTypeChat:
			client.processAndSaveChatMessage(&message)
		}

		client.Hub.Broadcast <- &message
	}
}

func (client *Client) WritePump() error {
	ticker := time.NewTicker(client.WebsocketSetting.PingPeriod)
	defer func() {
		ticker.Stop()
		client.CloseConnection()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Mu.Lock()
			client.Conn.SetWriteDeadline(time.Now().Add(client.WebsocketSetting.WriteTimeout))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed by server"))
				client.Mu.Unlock()
				return nil
			}

			writer, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				client.Mu.Unlock()
				return err
			}
			writer.Write(message)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				writer.Write(bytes.TrimSpace([]byte{'\n'}))
				writer.Write(<-client.Send)
			}
			writer.Close()
			client.Mu.Unlock()

		case <-ticker.C:
			client.Mu.Lock()
			client.Conn.SetWriteDeadline(time.Now().Add(client.WebsocketSetting.WriteTimeout))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				client.Mu.Unlock()
				return err
			}
			client.Mu.Unlock()

		case <-client.Done:
			return nil
		}
	}
}

func (client *Client) CloseConnection() {
	client.CloseOnce.Do(func() {
		close(client.Done)
		close(client.Send)
		client.Conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing connection"))
		client.Conn.Close()
	})
}

func (client *Client) IsClosed() bool {
	select {
	case <-client.Done:
		return true
	default:
		return false
	}
}

func (client *Client) processAndSaveChatMessage(message *Message) {
	var content string
	if err := json.Unmarshal(message.Content, &content); err != nil {
		return
	}
	savedMessage := client.ChatService.SaveMessage(client.RoomID, client.UserID, content)
	message.MessageID = savedMessage.ID
	message.SenderID = savedMessage.SenderID
}
