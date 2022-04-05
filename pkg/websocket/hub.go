package websocket

import (
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	ID             int    `json:"id,omitempty"`
	SenderUsername string `json:"sender,omitempty"`
	ThreadId       int    `json:"threadId,omitempty"`
	Body           string `json:"body,omitempty"`
}

type MessageChannel chan *Message
type UserChannel chan *Client

type Channels struct {
	messageChannel MessageChannel
	leaveChannel   UserChannel
	joinChannel    UserChannel
}

type Hub struct {
	clients  map[int][]*Client
	channels *Channels
}

func InitializeNewHub() *Hub {
	return &Hub{
		clients: make(map[int][]*Client),
		channels: &Channels{
			messageChannel: make(MessageChannel),
			leaveChannel:   make(UserChannel),
			joinChannel:    make(UserChannel),
		},
	}
}

func (h *Hub) addClient(c *Client) {
	/*
		 if userArray, ok := h.users[u.ThreadId]; ok { //checks if there already axists a thread with that id
				contains, i := Contains(userArray, u) //checks if the passed user already is in the thread
					if contains {
						userArray[i].Connection = u.Connection //readjust connection   ??maybe remove the userconnection pointer from the map and add the new
						log.Printf("reconnection") //do I even need this???
					} else {
						array := h.users[u.ThreadId]
						h.users[u.ThreadId] = append(array, u)
						log.Printf("User %s was added to thread %d.\n", u.Username, u.ThreadId)
					}
		} else { */
	h.clients[c.ThreadId] = append(h.clients[c.ThreadId], c) //adds the client to the array of clients (supports multiple connections to the same thread by the same user)
	fmt.Printf("User %s was added to thread %d.\n", c.Username, c.ThreadId)
	//}
}

func (h *Hub) removeClient(c *Client) {
	if userArray, ok := h.clients[c.ThreadId]; ok {
		i := Find(userArray, c) //finds the index (in the []*Client array) of the client that we want to close
		if i == -1 {
			fmt.Println("Error: Cannot kill connection")
			return
		}
		defer userArray[i].Connection.Close()
		h.clients[c.ThreadId] = RemoveIndex(userArray, i) //check if append creates new slice
		fmt.Printf("User %s left the chat\n", userArray[i].Username)
	}
}

func (h *Hub) SetupEventRouter() {
	for {
		select {
		case userjoin := <-h.channels.joinChannel:
			h.addClient(userjoin)
		case message := <-h.channels.messageChannel:
			for _, client := range h.clients[message.ThreadId] {
				select {
				case client.send <- message:
				default:
					close(client.send)
					h.removeClient(client)
				}
			}
		case userleave := <-h.channels.leaveChannel:
			h.removeClient(userleave)
		}
	}
}

func (h *Hub) ServeWs (w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{Connection: conn, send: make(chan *Message, 2000)}
	h.channels.joinChannel <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeStream(h.channels)
	go client.readStream(h.channels)
}