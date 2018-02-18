package server

import (
	"errors"

	"github.com/schollz/jsonstore"
)

var db *jsonstore.JSONStore

const databaseName = "textmailmachine.json.gz"

func popMessage(name string) (message messageJSON, err error) {
	var messages []messageJSON
	err = db.Get(name, &messages)
	if err != nil {
		return
	}
	if len(messages) > 0 {
		if len(messages) > 1 {
			db.Set(name, messages[1:])
		} else {
			db.Delete(name)
		}
		go jsonstore.Save(db, databaseName)
		message = messages[0]
	} else {
		err = errors.New("no messages")
	}
	return
}

func firstMessage(name string) (message messageJSON, err error) {
	var messages []messageJSON
	err = db.Get(name, &messages)
	if err != nil {
		return
	}
	if len(messages) > 0 {
		message = messages[0]
	} else {
		err = errors.New("no messages")
	}
	return
}

func addToDB(m messageJSON) (err error) {
	var messages []messageJSON
	err = db.Get(m.To, &messages)
	if err != nil {
		messages = []messageJSON{}
	}
	messages = append(messages, m)
	err = db.Set(m.To, messages)
	if err == nil {
		go jsonstore.Save(db, databaseName)
	}
	return
}

func getMessages(name string) (messages []messageJSON, err error) {
	err = db.Get(name, &messages)
	return
}
