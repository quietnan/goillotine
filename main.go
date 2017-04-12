package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

var someTestVideo = "JxWfvtnHtS0"

var db *databaseHandler

type audioRecord struct {
	File      string
	YoutubeID string
	Title     string
	Date      time.Time
}

func getAudio(youtubeID string) (*audioRecord, error) {
	for quality := -1; quality > -6; quality-- {
		fmt.Println("Trying quality ", quality)
		videoTitle, err := youtubeGet(youtubeID, quality) // "lala", error(nil)
		if err != nil {
			return nil, err
		}
		fmt.Println(videoTitle)
		outname := fmt.Sprintf("audioStore/%s_%v.mp3", youtubeID, quality)
		if err := transcode(youtubeID, outname, 64e3); err != nil {
			log.Printf("Could not extract audio: %v. Trying next better quality.", err)
		} else {
			log.Println("Success")
			return &audioRecord{File: outname, YoutubeID: youtubeID, Title: videoTitle, Date: time.Now()}, nil
		}
	}
	return nil, errors.New("Could not transcode any of the qualities")
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	youtubeID := vars["youtubeID"]
	fmt.Println("adding ", youtubeID)
	go func() {
		audio, err := getAudio(youtubeID)
		if err != nil {
			log.Fatalln("Could not get audio: ", err)
		}
		db.save(audio)
	}()
	w.Write([]byte("Scheduled " + youtubeID + " for download"))
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	audioRecords, err := db.listSince(time.Now())
	if err != nil {
		log.Fatalln("Could not get audio list: ", err)
	}

	feed := &feeds.Feed{
		Title: "My private feed",
		Link: &feeds.Link{
			Href: "http://localhost:8000/rss",
		},
		Created: time.Now(),
	}

	for _, ar := range audioRecords {
		feed.Add(&feeds.Item{
			Title: ar.Title,
			Link: &feeds.Link{
				Href:   "http://localhost:8000/" + ar.File,
				Type:   "audio/mpeg",
				Length: "1",
			},
			Created: ar.Date,
		})
	}
	rss, err := feed.ToRss()
	if err != nil {
		log.Fatalln("malformed feed", err)
	}
	w.Write([]byte(rss))
}

func main() {

	if _, err := os.Stat("./audioStore"); os.IsNotExist(err) {
		os.Mkdir("./audioStore", 0700)
	}

	var err error
	db, err = getDatabaseHandler("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer db.close()

	r := mux.NewRouter()
	r.HandleFunc("/add/{youtubeID}", addHandler)
	r.HandleFunc("/rss", feedHandler)
	r.PathPrefix("/audioStore/").Handler(http.StripPrefix("/audioStore/", http.FileServer(http.Dir("./audioStore"))))
	fmt.Println("listening...")
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatalln(err)
	}
}
