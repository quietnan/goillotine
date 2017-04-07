package main

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
)

type db struct {
	savechan  chan<- *audioRecord
	closechan chan bool
}

func getDatabaseHandler(address string) (*db, error) {
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
	return &db{savechan: savechan, closechan: closechan}, nil
}

func (self *db) save(content *audioRecord) {
	self.savechan <- content
}

func (self *db) close() {
	self.closechan <- true
	<-self.closechan
}
