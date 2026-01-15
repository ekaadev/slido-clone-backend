package websocket

import (
	"net/http"
	"slido-clone-backend/internal/util"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// for development
		return true
	},
}

type WebSocketHandler struct {
	hub          *Hub
	log          *logrus.Logger
	tokenUtil    *util.TokenUtil
	eventHandler *EventHandler
}

func NewWebSocketHandler(hub *Hub, log *logrus.Logger, tokenUtil *util.TokenUtil, eventHandler *EventHandler) *WebSocketHandler {
	return &WebSocketHandler{
		hub:          hub,
		log:          log,
		tokenUtil:    tokenUtil,
		eventHandler: eventHandler,
	}
}

// HandleWebSocket menangani koneksi WebSocket baru http -> ws
func (wsh *WebSocketHandler) HandleWebSocket(ctx *fiber.Ctx) error {
	// extract jwt token from query parameter
	token := ctx.Query("token")
	if token == "" {
		wsh.log.Warn("missing token parameter")
		return fiber.ErrBadRequest
	}

	// parse token and validate
	claims, err := wsh.tokenUtil.ParseToken(ctx.UserContext(), token)
	if err != nil {
		wsh.log.Warnf("invalid token: %v", err)
		return fiber.ErrUnauthorized
	}

	// validate required fields
	if claims.RoomID == nil || claims.ParticipantID == nil {
		wsh.log.Warn("missing required field in token")
		return fiber.ErrBadRequest
	}

	// upgrade http connection to websocket
	return adaptor.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			wsh.log.Error("upgrade failed: ", err)
			return
		}

		client := &Client{
			hub:            wsh.hub,
			conn:           conn,
			send:           make(chan []byte, 256),
			userID:         getUintValue(claims.UserID),
			roomID:         getUintValue(claims.RoomID),
			participantID:  getUintValue(claims.ParticipantID),
			isAnonymous:    claims.IsAnonymous,
			messageHandler: wsh.eventHandler.HandleMessage,
		}

		wsh.hub.register <- client

		go client.WritePump()
		go client.ReadPump()
	}))(ctx)
}

func getUintValue(ptr *uint) uint {
	if ptr == nil {
		return 0
	}
	return *ptr
}
