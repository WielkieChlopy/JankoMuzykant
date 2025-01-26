package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/zackradisic/soundcloud-api"
)

func main() {
	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{})

	if err != nil {
		log.Fatal(err.Error())
	}

	track, err := sc.GetTrackInfo(soundcloudapi.GetTrackInfoOptions{
		URL: "https://on.soundcloud.com/wqWaV11nwvDi9xHM8",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	jsonm, err := json.Marshal(track)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(string(jsonm))
}
