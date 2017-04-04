package main

import (
	"errors"
	"fmt"
	"log"
)

var someTestVideo = "JxWfvtnHtS0"

func getAudio(youtubeId string) (string, error) {
	for quality := -1; quality > -6; quality-- {
		fmt.Println("Trying quality ", quality)
		if videoDesc, err := youtubeGet(youtubeId, quality); err != nil {
			return "", err
		} else {
			fmt.Println(videoDesc)
		}
		outname := fmt.Sprintf("%s_%v.mp3", youtubeId, quality)
		if err := transcode(youtubeId, outname, 23e3); err != nil {
			log.Printf("Could not extract audio: %v. Trying next better quality.", err)
		} else {
			log.Println("Success")
			return outname, nil
		}
	}
	return "", errors.New("Could not transcode any of the qualities")
}

func main() {
	if file, err := getAudio(someTestVideo); err != nil {
		log.Fatal("Could not get audio: ", err)
	} else {
		fmt.Println("The audio is now at: ", file)
	}
}
