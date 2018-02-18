package server

import (
	"encoding/json"
	"fmt"
	"sync"

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

	// time to next message
	Queue messageQueue

	hasMessage bool
}

type messageQueue struct {
	Messages []messageJSON
	sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		hasMessage: false,
		Queue:      messageQueue{Messages: []messageJSON{}},
	}
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

func (h *Hub) handleMessage(m messageJSON) (err error) {
	m, err = validateMessage(m)
	if err != nil {
		return
	}
	h.Queue.Lock()
	h.Queue.Messages = append(h.Queue.Messages, m)
	h.Queue.Unlock()
	go h.broadcastNextMessage(false)
	return
}

func (h *Hub) broadcastNextMessage(force bool) {
	// overwrite current message only if forced
	// or if there is currently no message
	if !force && h.hasMessage {
		log.Debug("not sending out message")
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
		messageHTML.From = fmt.Sprintf("- %s<br>(%s)<br>Seen by %d.", message.From, humanize.Time(message.Timestamp), len(h.clients))
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
