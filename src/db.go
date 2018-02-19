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

func (d *Database) close() {
	err := d.fileLock.Unlock()
	if err != nil {
		log.Error(err)
	} else {
		os.Remove(path.Join("data", d.Name+".lock"))
	}
}

func (d *Database) saveMessage(message messageJSON) (err error) {
	log.Debug("saving message")
	var messages []messageJSON
	messagesB, errRead := ioutil.ReadFile(path.Join("data", d.Name+".json"))
	log.Debug("loaded messages")
	if errRead == nil {
		err = json.Unmarshal(messagesB, &messages)
		if err != nil {
			return err
		}
		log.Debug("unmarshaled")
	} else {
		messages = []messageJSON{}
	}
	messages = append(messages, message)
	log.Debugf("have %d messages for %s", len(messages), d.Name)
	err = d.saveMessages(messages)
	log.Debugf("saved messages for %s", d.Name)
	return
}

func (d *Database) saveMessages(messages []messageJSON) (err error) {
	if len(messages) == 0 {
		log.Debug("removing database")
		err = os.Remove(path.Join("data", d.Name+".json"))
	} else {
		messageBytes, errMarshal := json.Marshal(messages)
		if errMarshal != nil {
			return errMarshal
		}
		err = ioutil.WriteFile(path.Join("data", d.Name+".json"), messageBytes, 0644)
		log.Debugf("wrote %d messages", len(messages))
	}
	return
}

func (d *Database) popMessage() (messages []messageJSON, err error) {
	messagesB, errRead := ioutil.ReadFile(path.Join("data", d.Name+".json"))
	if errRead != nil {
		err = errors.New("no messages")
		return
	}
	err = json.Unmarshal(messagesB, &messages)
	if err != nil {
		return
	}
	if len(messages) == 1 {
		err = d.saveMessages([]messageJSON{})
	} else {
		err = d.saveMessages(messages[1:])
	}
	return
}
