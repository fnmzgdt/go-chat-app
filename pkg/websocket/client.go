package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)
/*
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)
*/
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	Username   string
	ThreadId   int
	send 	   chan *Message
	Connection *websocket.Conn
}

func (c *Client) readStream(channels *Channels) { //pass ws.channel
	defer func() {
		channels.leaveChannel <- c 	//push the user chat into the 1st user channel (LEAVE)
		c.Connection.Close()
	}()
	c.Connection.SetReadLimit(maxMessageSize)
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.SetPongHandler(func(string) error { c.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, messageJson, err := c.Connection.ReadMessage() //if theres a message sent from the client get it
		if err != nil {
			log.Println(err.Error())
			break
		}

		message := &Message{}
		if err := json.Unmarshal(messageJson, message); err != nil {
			log.Printf("json error: %v", err)
		}

		channels.messageChannel <- message //push the message struct into the message channel
	}
}

func (c *Client) writeStream(channels *Channels) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
			}

			w, err := c.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			messageJSON, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}

			w.Write(messageJSON)

			n := len(c.send) // Check for queued chat messages and if there are any add them to the current websocket message.
			for i := 0; i < n; i++ {
				messageJSON, err = json.Marshal(<-c.send)
				if err != nil {
					fmt.Printf("Error: %v", err)
				}

				w.Write(messageJSON)
			}
			
			if err := w.Close(); err != nil {
				return
			}

			case <-ticker.C:
				c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
		}
	}
}