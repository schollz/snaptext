package server

import (
	"errors"
	"fmt"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
)

type messageHTML struct {
	Meta       string `json:"meta"`
	Message    string `json:"message"`
	Submessage string `json:"submessage"`
}

type messageJSON struct {
	To        string    `json:"to" binding:"required"`
	From      string    `json:"from" binding:"required"`
	Message   string    `json:"message" binding:"required"`
	Display   int       `json:"display"`
	Timestamp time.Time `json:"timestamp"`
}

func validateMessage(m messageJSON) (messageJSON, error) {
	var err error
	if m.Display == 0 {
		m.Display = 10
	}
	m.Timestamp = time.Now()
	m.To = strings.ToLower(strings.TrimSpace(m.To))
	m.From = strings.ToLower(strings.TrimSpace(m.From))
	m.Message = strings.TrimSpace(m.Message)
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowAttrs("href").OnElements("a")
	p.AllowElements("p")
	p.AllowElements("em")
	p.AllowElements("i")
	p.AllowElements("small")
	m.Message = p.Sanitize(m.Message)
	if len(m.To) == 0 || len(m.From) == 0 || len(m.Message) == 0 {
		err = errors.New("to, from, and message cannot be empty")
	}
	return m, err
}

func getNextMessage(name string) (m messageHTML, err error) {
	db := open(name)
	defer db.close()

	messages, err := db.popMessage()

	if err != nil {
		m.Message = "No messages."
	} else {
		m.Message = messages[0].Message
		m.Submessage = fmt.Sprintf("Sent from <a class='link dim mid-gray' href='/?to=%s&from=%s'>%s</a> %s.", strings.ToLower(messages[0].From), name, messages[0].From, humanize.Time(messages[0].Timestamp))
		if len(messages) > 1 {
			m.Meta = "more messages"
		}
	}
	return
}
