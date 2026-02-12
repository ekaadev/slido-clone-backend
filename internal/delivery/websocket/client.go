package websocket

import (
	"time"

	"github.com/gofiber/contrib/websocket"
)

const (
	// waktu maksimal untuk menulis pesan ke client dalam detik
	writeWait = 10 * time.Second

	// waktu maksimal untuk menunggu pong dari client dalam detik
	pongWait = 60 * time.Second

	// periode untuk mengirim ping ke client dalam detik, pingPeriod < pongWait
	pingPeriod = (pongWait * 9) / 10

	// ukuran maksimal pesan dalam byte
	maxMessage = 512 * 1024
)

// Client representasi koneksi websocket ke client
type Client struct {
	hub  *Hub            // hub websocket
	conn *websocket.Conn // koneksi websocket
	send chan []byte     // channel untuk mengirim pesan ke client

	// identitas client
	userID        uint // dari jwt token
	roomID        uint // room yang dijoin
	participantID uint // dari participant
	isAnonymous   bool // anonymous
	isRoomOwner   bool // true jika user adalah pembuat room (host)

	// handler reference (untuk process events)
	messageHandler func(*Client, []byte) error
}

// ReadPump goroutine untuk membaca pesan dari client
func (c *Client) ReadPump() {
	defer func() {
		if r := recover(); r != nil {
			c.hub.log.Warnf("websocket read pump panic: %v", r)
		}
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// set configurasi koneksi
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// loop membaca pesan
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.hub.log.WithField("error", err).Error("WebSocket read error")
			}
			break
		}

		// untuk kebutuhan testing kita echo kembali pesan yang diterima
		c.hub.log.WithField("message", msg).Debug("WebSocket message received")

		if c.messageHandler != nil {
			if err = c.messageHandler(c, msg); err != nil {
				c.hub.log.Warnf("failed to handle message: %+v", err)
			}
		}
	}
}

// WritePump goroutine untuk menulis pesan ke client
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if r := recover(); r != nil {
			c.hub.log.Warnf("websocket write pump panic: %v", r)
		}
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			// kirim semua pending message ke queue
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
