package websocket

import (
	"slido-clone-backend/internal/util"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

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

	// upgrade to websocket
	// call two functions: first is to create the websocket connection handler,
	// second is the actual handler returned by websocket.New()
	return websocket.New(func(c *websocket.Conn) {
		wsh.log.Debug("upgraded to websocket connection")

		defer func() {
			if r := recover(); r != nil {
				wsh.log.Warnf("websocket connection handler panic: %v", r)
			}
		}()

		// create new client
		client := &Client{
			hub:            wsh.hub,
			conn:           c,
			send:           make(chan []byte, 256),
			userID:         getUintValue(claims.UserID),
			roomID:         getUintValue(claims.RoomID),
			participantID:  getUintValue(claims.ParticipantID),
			displayName:    claims.DisplayName,
			isAnonymous:    claims.IsAnonymous,
			isRoomOwner:    claims.IsRoomOwner,
			messageHandler: wsh.eventHandler.HandleMessage,
		}

		// register client ke hub
		client.hub.register <- client

		// run write pump sebagai goroutine
		go client.WritePump()

		// run read pump di goroutine utama agar koneksi tetap terbuka
		client.ReadPump()

		// cleanup after disconnect
		wsh.eventHandler.HandleDisconnect(client)
	})(ctx)
}

func getUintValue(ptr *uint) uint {
	if ptr == nil {
		return 0
	}
	return *ptr
}
