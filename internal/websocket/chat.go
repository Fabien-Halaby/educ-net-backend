package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type Client struct {
	ID   int
	Conn *websocket.Conn
	Send chan Message
}

type Room struct {
	ID         int
	Clients    map[*Client]bool
	mu         sync.RWMutex
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

var rooms = make(map[int]*Room)

func GetRoom(classID int) *Room {
	room, exists := rooms[classID]
	if !exists {
		room = &Room{
			ID:         classID,
			Clients:    make(map[*Client]bool),
			Broadcast:  make(chan Message),
			Register:   make(chan *Client),
			Unregister: make(chan *Client),
		}
		rooms[classID] = room
		go room.Run()
	}
	return room
}

func NewRoom(classID int) *Room {
	return GetRoom(classID)
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.mu.Lock()
			r.Clients[client] = true
			r.mu.Unlock()
			log.Printf("Client %d joined room %d (%d total)", client.ID, r.ID, len(r.Clients))

		case client := <-r.Unregister:
			r.mu.Lock()
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}
			r.mu.Unlock()
			log.Printf("Client %d left room %d (%d total)", client.ID, r.ID, len(r.Clients))

		case message := <-r.Broadcast:
			r.mu.RLock()
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					r.mu.RUnlock()
					r.Unregister <- client
					r.mu.RLock()
				}
			}
			r.mu.RUnlock()
		}
	}
}

func (r *Room) RegisterClient(client *Client) {
	r.Register <- client
}

func (r *Room) UnregisterClient(client *Client) {
	r.Unregister <- client
}
