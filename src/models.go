package server

import (
	"errors"
	"strings"
	"time"

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
