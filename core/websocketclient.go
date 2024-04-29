// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}

// wsReader pumps messages from the websocket connection to the hub.
//
// The application runs wsReader in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) wsReader() {
	var message Message

	// unregister client if channel is closed
	defer func() {
		logrus.WithFields(logrus.Fields{
			"remote_address": c.Conn.RemoteAddr(),
		}).Debug("websocket client unregistered and connection closed")
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// set connection specifics
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// read channel
	for {
		// read message
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			// handle connection errors
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithFields(logrus.Fields{
					"error":          err.Error(),
					"remote_address": c.Conn.RemoteAddr(),
				}).Error("websocket connection closed from remote client")
			}

			// channel closed, cannot send anything to client
			break
		}

		// parse message to a go struct
		err = json.Unmarshal(bytes.TrimSpace(bytes.Replace(msg, newline, space, -1)), &message)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":       err.Error(),
				"raw_message": string(msg),
			}).Error("invalid json from remote websocket client")

			// generate response for client, always send something to client
			c.Send <- NewInvalidFormatErrorMessage().JSONMarshal()
			continue
		}

		// override timestamp to avoid any kind of issue
		message.Timestamp = time.Now()
		// also check for valid UUID or generate a new one
		if _, err := uuid.ParseBytes([]byte(message.ID)); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err.Error(),
				"message_id":     message.ID,
				"remote_address": c.Conn.RemoteAddr(),
			}).Debug("error parsing message uuid from client")
			message.ID = uuid.New().String()
		}

		// inspect message type
		switch message.Type {
		case "chat":
			// chat message must be broadcasted to all clients
			msg, err := json.Marshal(message)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"remote_address": c.Conn.RemoteAddr(),
					"error":          err.Error(),
					"raw_message":    message,
				}).Error("error generating json response for client")

				// generate response for client, always send something to client
				c.Send <- NewMarshalJSONErrorMessage().JSONMarshal()
			}

			// send message to broadcast channel
			c.Hub.Broadcast <- msg
		default:
			// unrecognized message type
			logrus.WithFields(logrus.Fields{
				"remote_address": c.Conn.RemoteAddr(),
				"raw_message":    message,
			}).Error("invalid message type in remote websocket message")

			// generate response for client, always send something to client
			c.Send <- NewIvalidTypeErrorMessage().JSONMarshal()
		}
	}
}

// wsWriter pumps messages from the hub to the websocket connection.
//
// A goroutine running wsWriter is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) wsWriter() {
	// start ticker to send ping to clients
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			// set write timeout on connection
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// get a new writer and close the previous one if still opened
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// generate message for client, always send something to client
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			// close writer when done
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// send ping message on ticker emit
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			// close channel on error
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
