package main

import (
	"errors"
	"fmt"
	"log"
)

var someTestVideo = "JxWfvtnHtS0"

type audioRecord struct {
	File      string
	YoutubeID string
	Title     string
}

func getAudio(youtubeID string) (*audioRecord, error) {
	for quality := -1; quality > -6; quality-- {
		fmt.Println("Trying quality ", quality)
		videoTitle, err := youtubeGet(youtubeID, quality)
		if err != nil {
			return nil, err
		}
		fmt.Println(videoTitle)
		outname := fmt.Sprintf("%s_%v.mp3", youtubeID, quality)
		if err := transcode(youtubeID, outname, 23e3); err != nil {
			log.Printf("Could not extract audio: %v. Trying next better quality.", err)
		} else {
			log.Println("Success")
			return &audioRecord{File: outname, YoutubeID: youtubeID, Title: videoTitle}, nil
		}
	}
	return nil, errors.New("Could not transcode any of the qualities")
}

func main() {
	newAudio, close := getDatabaseHandlers()
	defer func() { close <- true }()

	audio, err := getAudio(someTestVideo)

	if err != nil {
		log.Fatalln("Could not get audio: ", err)
	}
	newAudio <- audio
}
