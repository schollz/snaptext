package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
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
	m.To = strings.TrimSpace(m.To)
	m.From = strings.TrimSpace(m.From)
	m.Message = strings.TrimSpace(m.Message)
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowAttrs("href").OnElements("a")
	p.AllowElements("p")
	p.AllowElements("em")
	p.AllowElements("i")
	m.Message = p.Sanitize(m.Message)
	if len(m.To) == 0 || len(m.From) == 0 || len(m.Message) == 0 {
		err = errors.New("to, from, and message cannot be empty")
	}
	return m, err
}

func saveMessages(name string, messages []messageJSON) (err error) {
	messageQueue, err := json.Marshal(h.Queue.Messages)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path.Join("data", name+".json"), messageQueue, 0644)
	return
}

func saveMessage(name string, message messageJSON) (err error) {
	var messages []messageJSON
	messagesB, errRead := ioutil.ReadFile(path.Join("data", name+".json"))
	if errRead == nil {
		err = json.Unmarshal(messagesB, &messages)
		if err == nil {
			return err
		}
	} else {
		messages = []messageJSON{}
	}
	messages = append(messages, message)
	messageQueue, err := json.Marshal(messages)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path.Join("data", name+".json"), messageQueue, 0644)
	return
}

func popMessage(name string) (message messageJSON, err error) {
	messagesB, errRead := ioutil.ReadFile(path.Join("data", name+".json"))
	if errRead != nil {
		return errors.New("no messages")
	}
	var messages []messageJSON
	err = json.Unmarshal(messagesB, &messages)
	if err == nil {
		return err
	}
	message = messages[0]
	if len(messages) == 1 {
		messages = []messageJSON{}
	} else {
		messages = messages[1:]
	}
	err = saveMessages(name, messages)
	return
}
