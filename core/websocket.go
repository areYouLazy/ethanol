package core

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	webSocketHub *Hub

	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  viper.GetInt("ethanol.server.websocket.read_buffer_size"),
	WriteBufferSize: viper.GetInt("ethanol.server.websocket.write_buffer_size"),
	CheckOrigin: func(r *http.Request) bool {
		header := r.Header.Get("Origin")
		logrus.WithFields(logrus.Fields{
			"origin_header_value": header,
			"remote_ip":           r.RemoteAddr,
		}).Debug("analyzing origin header to allow websocket upgrade request")

		for _, origin := range viper.GetStringSlice("ethanol.server.websocket.origins") {
			if header == origin {
				return true
			}
		}

		return false
	},
}

//InitWebSocketHub initialize websocket hub
func InitWebSocketHub() *Hub {
	u := newHub()
	go u.Run()

	return u
}

// serveWebSocket handles websocket requests from the peer.
func serveWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":          err.Error(),
			"remote_address": r.RemoteAddr,
		}).Error("error upgrading connection to websocket")
		return
	}

	// instantiate new client with reference to ws connection and hub
	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	// register client to hub
	client.Hub.Register <- client
	logrus.WithFields(logrus.Fields{
		"remote_address": client.Conn.RemoteAddr(),
	}).Debug("new websocket client registered")

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.wsWriter()
	go client.wsReader()
}
