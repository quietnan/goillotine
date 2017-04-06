package main

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
)

func getDatabaseHandlers() (chan<- *audioRecord, chan<- bool) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalln(err)
	}
	c := session.DB("goillotine").C("audio")
	save := make(chan *audioRecord)
	close := make(chan bool)
	go func() {
		for {
			select {
			case a := <-save:
				if err := c.Insert(a); err != nil {
					log.Fatalln(err)
				}
				fmt.Println("added: ", a)
			case <-close:
				session.Close()
				break
			}
		}
	}()
	return save, close
}
