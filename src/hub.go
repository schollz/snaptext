package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/cihub/seelog"
	humanize "github.com/dustin/go-humanize"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	Name string

	hasMessage bool
}

func newHub(name string) *Hub {
	h := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		hasMessage: false,
		Name:       name,
	}
	return h
}

func (h *Hub) run() {
	log.Debug("starting new hub")
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				if len(h.clients) == 0 {
					h.hasMessage = false
				}
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

// serveWs handles websocket requests from the peer.
func (h *Hub) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func (h *Hub) broadcastNextMessage(force bool) {
	defer h.saveMessages()
	// overwrite current message only if forced
	// or if there is currently no message
	if !force && h.hasMessage {
		log.Debug("not sending out message")
		h.broadcast <- []byte(`{"meta":"new"}`)
		return
	}
	h.Queue.Lock()
	var messageHTML messageHTML
	if len(h.Queue.Messages) == 0 {
		messageHTML.Message = "No messages."
		h.hasMessage = false
	} else {
		message := h.Queue.Messages[0]
		if len(h.Queue.Messages) == 1 {
			h.Queue.Messages = []messageJSON{}
		} else {
			h.Queue.Messages = h.Queue.Messages[1:]
		}
		messageHTML.Message = message.Message
		messageHTML.Submessage = fmt.Sprintf("Sent from %s %s.", message.From, humanize.Time(message.Timestamp))
		if len(h.Queue.Messages) > 0 {
			messageHTML.Meta = "more messages"
		}
		h.hasMessage = true
	}
	h.Queue.Unlock()

	bMessage, errMarshal := json.Marshal(messageHTML)
	if errMarshal != nil {
		log.Warn(errMarshal)
		return
	}
	h.broadcast <- bMessage
}