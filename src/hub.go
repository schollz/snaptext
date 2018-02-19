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
	db := open(h.Name)
	defer db.close()

	// overwrite current message only if forced
	// or if there is currently no message
	if !force && h.hasMessage {
		log.Debug("not sending out message")
		h.broadcast <- []byte(`{"meta":"new"}`)
		return
	}

	messages, err := db.popMessage()

	var messageHTML messageHTML
	if err != nil {
		messageHTML.Message = "No messages."
		h.hasMessage = false
	} else {
		messageHTML.Message = messages[0].Message
		messageHTML.Submessage = fmt.Sprintf("Sent from %s %s.", messages[0].From, humanize.Time(messages[0].Timestamp))
		if len(messages) > 0 {
			messageHTML.Meta = "more messages"
		}
		h.hasMessage = true
	}

	bMessage, errMarshal := json.Marshal(messageHTML)
	if errMarshal != nil {
		log.Warn(errMarshal)
		return
	}
	h.broadcast <- bMessage
}
