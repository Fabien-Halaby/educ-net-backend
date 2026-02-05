package handler

import (
	"context"
	"educnet/internal/middleware"
	"educnet/internal/usecase"
	"educnet/internal/utils"
	ws "educnet/internal/websocket"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	uc usecase.MessageUseCase
}

func NewChatHandler(uc usecase.MessageUseCase) *ChatHandler {
	return &ChatHandler{uc}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *ChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("🔥 WS Handler START - classId parsing...")

	classIDStr := mux.Vars(r)["classId"]
	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		utils.BadRequest(w, "Invalid class id type")
		return
	}
	log.Printf("✅ WS classID parsed: %d", classID)

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.Unauthorized(w, "Unauthorized user")
		return
	}
	log.Printf("✅ WS UserID: %d (%s)", claims.UserID, claims.Email)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade:", err)
		return
	}
	defer conn.Close()

	canAccess, err := h.uc.CanAccessClass(r.Context(), claims.UserID, classID)
	if err != nil {
		utils.HandleUseCaseError(w, err)
		return
	}
	if !canAccess {
		conn.WriteJSON(map[string]string{"error": "Access denied"})
		return
	}

	room := ws.NewRoom(classID)
	client := &ws.Client{
		ID:   claims.UserID,
		Conn: conn,
		Send: make(chan ws.Message, 256),
	}

	room.RegisterClient(client)

	messages, err := h.uc.GetClassMessages(r.Context(), classID, 50)
	if err == nil {
		for _, msg := range messages {
			roomMsg := ws.Message{
				Type:    "message",
				Content: msg,
			}
			select {
			case client.Send <- roomMsg:
			default:
				//! Client lent
			}
		}
	}

	go h.writePump(conn, client.Send)

	h.readPump(room, conn, client, claims.UserID, classID)
}

func (h *ChatHandler) writePump(conn *websocket.Conn, send chan ws.Message) {
	ticker := time.NewTicker(90 * time.Second)
	defer ticker.Stop()
	defer conn.Close()

	for {
		select {
		case message, ok := <-send:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteJSON(message); err != nil {
				log.Println("Write error:", err)
				return
			}

		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *ChatHandler) readPump(room *ws.Room, conn *websocket.Conn, client *ws.Client, userID, classID int) {
	defer func() {
		room.UnregisterClient(client)
		conn.Close()
	}()

	for {
		var msg struct {
			Content string `json:"content"`
		}

		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			break
		}

		if msg.Content == "" || len(msg.Content) > 1000 {
			continue
		}

		createdMsg, err := h.uc.SendMessage(context.Background(), userID, classID, msg.Content)
		if err != nil {
			conn.WriteJSON(map[string]string{"error": err.Error()})
			continue
		}

		roomMsg := ws.Message{
			Type:    "message",
			Content: createdMsg,
		}
		room.Broadcast <- roomMsg
	}
}
