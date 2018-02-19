package server

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"

	log "github.com/cihub/seelog"
	flock "github.com/theckman/go-flock"
)

type Database struct {
	Name     string
	fileLock *flock.Flock
}

func open(name string) *Database {
	d := new(Database)
	d.Name = name
	d.fileLock = flock.NewFlock(path.Join("data", name+".lock"))
	for {
		locked, err := d.fileLock.TryLock()
		if err == nil && locked {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return d
}

func (d *Database) unlock() {
	err := d.fileLock.Unlock()
	if err != nil {
		log.Error(err)
	} else {
		os.Remove(path.Join("data", d.Name+".lock"))
	}
}

func (d *Database) saveMessages(messages []messageJSON) (err error) {
	messageBytes, err := json.Marshal(messages)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path.Join("data", d.Name+".json"), messageBytes, 0644)
	return
}

func (d *Database) saveMessage(message messageJSON) (err error) {
	var messages []messageJSON
	messagesB, errRead := ioutil.ReadFile(path.Join("data", d.Name+".json"))
	if errRead == nil {
		err = json.Unmarshal(messagesB, &messages)
		if err == nil {
			return err
		}
	} else {
		messages = []messageJSON{}
	}
	messages = append(messages, message)
	messageBytes, err := json.Marshal(messages)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path.Join("data", d.Name+".json"), messageBytes, 0644)
	return
}

func (d *Database) popMessage() (message messageJSON, err error) {
	messagesB, errRead := ioutil.ReadFile(path.Join("data", d.Name+".json"))
	if errRead != nil {
		err = errors.New("no messages")
		return
	}
	var messages []messageJSON
	err = json.Unmarshal(messagesB, &messages)
	if err != nil {
		return
	}
	message = messages[0]
	if len(messages) == 1 {
		messages = []messageJSON{}
	} else {
		messages = messages[1:]
	}
	err = saveMessages(d.Name, messages)
	return
}
