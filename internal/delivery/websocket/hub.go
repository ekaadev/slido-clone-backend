package websocket

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

type Hub struct {
	clients    map[*Client]bool          // registered clients
	broadcast  chan []byte               // kirim pesan ke semua client
	register   chan *Client              // client yang mau register
	unregister chan *Client              // client yang mau unregister
	rooms      map[uint]map[*Client]bool // rooms dan clients di dalamnya
	log        *logrus.Logger
}

// NewHub membuat instance Hub baru
func NewHub(log *logrus.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256), // buffered channel -> ukuran channel yang reasonable agar tidak memakan memori berlebihan
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[uint]map[*Client]bool),
		log:        log,
	}
}

// Run goroutine untuk mengelola semua channel operation
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			if h.rooms[client.roomID] == nil {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}

			h.rooms[client.roomID][client] = true

			h.log.WithFields(logrus.Fields{
				"user_id":        client.userID,
				"room_id":        client.roomID,
				"participant_id": client.participantID,
			}).Info("Client connected")

			// broadcast participant joined ke room
			h.broadcastParticipantJoined(client)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				// broadcast participant left ke room sebelum unregister
				h.broadcastParticipantLeft(client)

				delete(h.clients, client)
				close(client.send)

				if h.rooms[client.roomID] != nil {
					delete(h.rooms[client.roomID], client)

					if len(h.rooms[client.roomID]) == 0 {
						delete(h.rooms, client.roomID)
					}
				}

				h.log.WithFields(logrus.Fields{
					"user_id": client.userID,
					"room_id": client.roomID,
				}).Debug("Client disconnected")
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// BroadcastToRoom mengirim pesan ke semua client di room tertentu
func (h *Hub) BroadcastToRoom(roomID uint, msg []byte) {
	if clients, ok := h.rooms[roomID]; ok {
		h.log.WithFields(logrus.Fields{
			"room_id":      roomID,
			"client_count": len(clients),
		}).Debug("Broadcasting to room")

		for client := range clients {
			select {
			case client.send <- msg:
			default:
				h.log.WithFields(logrus.Fields{
					"user_id": client.userID,
				}).Debug("Failed to send message to client")
			}
		}
	}
}

// broadcastParticipantJoined broadcast ketika participant baru join room
func (h *Hub) broadcastParticipantJoined(client *Client) {
	data := WSMessage{
		Event: EventRoomUserJoin,
		Data: h.mustMarshal(map[string]interface{}{
			"participant_id": client.participantID,
			"display_name":   client.displayName,
			"is_anonymous":   client.isAnonymous,
			"joined_at":      time.Now().Format(time.RFC3339),
		}),
	}
	h.BroadcastToRoom(client.roomID, h.mustMarshal(data))
}

// broadcastParticipantLeft broadcast ketika participant disconnect dari room
func (h *Hub) broadcastParticipantLeft(client *Client) {
	data := WSMessage{
		Event: EventRoomUserLeft,
		Data: h.mustMarshal(map[string]interface{}{
			"participant_id": client.participantID,
			"display_name":   client.displayName,
			"left_at":        time.Now().Format(time.RFC3339),
		}),
	}
	h.BroadcastToRoom(client.roomID, h.mustMarshal(data))
}

// mustMarshal helper untuk marshal JSON, panic jika error
func (h *Hub) mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		h.log.WithField("error", err).Error("failed to marshal JSON")
		return []byte("{}")
	}
	return data
}
