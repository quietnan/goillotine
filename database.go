package main

import (
	"fmt"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type databaseHandler struct {
	savechan  chan<- *audioRecord
	closechan chan bool
	session   *mgo.Session
}

func getDatabaseHandler(address string) (*databaseHandler, error) {
	session, err := mgo.Dial(address)
	if err != nil {
		return nil, err
	}
	c := session.DB("goillotine").C("audio")
	savechan := make(chan *audioRecord)
	closechan := make(chan bool)
	go func() {
		for {
			select {
			case a := <-savechan:
				if err := c.Insert(a); err != nil {
					log.Fatalln(err)
				}
				fmt.Println("added: ", a)
			case <-closechan:
				session.Close()
				closechan <- true
				break
			}
		}
	}()
	return &databaseHandler{savechan: savechan, closechan: closechan, session: session}, nil
}

func (d *databaseHandler) save(content *audioRecord) {
	d.savechan <- content
}

func (d *databaseHandler) listSince(start time.Time) ([]*audioRecord, error) {
	sessionCopy := d.session.Copy()
	defer sessionCopy.Close()

	c := sessionCopy.DB("goillotine").C("audio")

	var ret []*audioRecord
	err := c.Find(nil).All(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (d *databaseHandler) close() {
	d.closechan <- true
	<-d.closechan
}
