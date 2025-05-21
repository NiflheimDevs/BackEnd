package websocket

import (
	"encoding/json"
	"fmt"
	"log"
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
				if client.IsClosed() {
					delete(room, client)
					continue
				}
				select {
				case client.Send <- messageBytes:
				default:
					go func(c *Client) {
						select {
						case hub.Unregister <- c:
						default:
							log.Println("Unregister channel full while broadcasting")
						}
					}(client)
				}
			}
		}
	case MessageTypeNotification:
		if clients, ok := hub.Clients[message.SenderID]; ok {
			for client := range clients {
				if client.IsClosed() {
					delete(clients, client)
					continue
				}
				select {
				case client.Send <- message.Content:
				default:
					go func(c *Client) {
						select {
						case hub.Unregister <- c:
						default:
							log.Println("Unregister channel full while sending notification")
						}
					}(client)
				}
			}
		}
	}
}

func (hub *Hub) SendToUser(userID int, messageType string, content []byte) error {
	hub.Mu.RLock()
	clients, ok := hub.Clients[userID]
	hub.Mu.RUnlock()

	if !ok || len(clients) == 0 {
		return fmt.Errorf("no active connections for user %d", userID)
	}

	for client := range clients {
		if client.IsClosed() {
			hub.Unregister <- client
			continue
		}
		select {
		case client.Send <- content:
		default:
			go func(c *Client) {
				select {
				case hub.Unregister <- c:
				default:
					log.Println("Unregister channel full while sending to user")
				}
			}(client)
		}
	}
	return nil
}
