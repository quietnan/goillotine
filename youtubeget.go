package main

import (
	"fmt"
	"log"

	youtube "github.com/KeluDiao/gotube/api"
)

func youtubeGet(target string, quality int) (string, error) {
	log.Println("download to file=", target)
	vl, err := youtube.GetVideoListFromId(target)
	if err != nil {
		fmt.Println("COULD NOT GET VIDEO LIST")
		return "", err
	}
	if quality < 0 {
		quality += len(vl.Videos)
	}
	fmt.Println("Downloading ", vl.Title, "in quality ", quality)
	v := (vl.Videos[quality])
	if err := v.Download("", target); err != nil {
		fmt.Println("COULD NOT DOWNLOAD")
		return "", err
	}
	return v.Title, nil
}
