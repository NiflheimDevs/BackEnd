package websocket

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	Clients    map[int]map[*Client]bool
	Rooms      map[int]map[*Client]bool
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]map[*Client]bool),
		Rooms:      make(map[int]map[*Client]bool),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.handleRegister(client)
		case client := <-hub.Unregister:
			hub.handleUnregister(client)
		case message := <-hub.Broadcast:
			hub.handleBroadcast(message)
		}
	}
}

func (hub *Hub) handleRegister(client *Client) {
	hub.Mu.Lock()
	defer hub.Mu.Unlock()

	if _, ok := hub.Clients[client.UserID]; !ok {
		hub.Clients[client.UserID] = make(map[*Client]bool)
	}
	hub.Clients[client.UserID][client] = true

	if _, ok := hub.Rooms[client.RoomID]; !ok {
		hub.Rooms[client.RoomID] = make(map[*Client]bool)
	}
	hub.Rooms[client.RoomID][client] = true
}

func (hub *Hub) handleUnregister(client *Client) {
	hub.Mu.Lock()
	defer hub.Mu.Unlock()

	if clients, ok := hub.Clients[client.UserID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(hub.Clients, client.UserID)
		}
	}

	if room, ok := hub.Rooms[client.RoomID]; ok {
		delete(room, client)
		if len(room) == 0 {
			delete(hub.Rooms, client.RoomID)
		}
	}

	client.CloseConnection()
}

func (hub *Hub) handleBroadcast(message *Message) {
	hub.Mu.Lock()
	defer hub.Mu.Unlock()

	switch message.Type {
	case MessageTypeChat:
		messageBytes, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}
		if room, ok := hub.Rooms[message.RoomID]; ok {
			for client := range room {
				select {
				case client.Send <- messageBytes:
				default:
					hub.Unregister <- client
				}
			}
		}
	case MessageTypeNotification:
		if clients, ok := hub.Clients[message.SenderID]; ok {
			for client := range clients {
				select {
				case client.Send <- message.Content:
				default:
					hub.Unregister <- client
				}
			}
		}
	}
}

func (hub *Hub) SendToUser(userID int, messageType string, content []byte) {
	hub.Mu.RLock()
	defer hub.Mu.RUnlock()

	if clients, ok := hub.Clients[userID]; ok {
		for client := range clients {
			select {
			case client.Send <- content:
			default:
				hub.Unregister <- client
			}
		}
	}
}
