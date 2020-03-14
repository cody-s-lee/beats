package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/benbjohnson/clock"
	"github.com/cody-s-lee/beats/beats"
)

func main() {
	args := os.Args[1:]

	var song beats.Song
	if len(args) > 0 {
		fn := args[0]
		reader, err := os.Open(fn)
		if err != nil {
			log.Fatalf("%s\nFilename provided: %s\n", err, fn)
		}

		song = getSong(reader)
	} else {
		song = getDefaultSong()
	}

	bytes, err := json.Marshal(song)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(bytes))

	fmt.Printf("Name: %s\n", song.Name)
	fmt.Printf("Tempo: %d bpm\n", song.Tempo)

	clock := clock.New()
	out := make(chan beats.Step)

	go song.Play(clock, out)

	for s := range out {
		fmt.Printf("%d: %s\n", s.Tick, s.Beat)
	}
}

func getSong(reader io.Reader) beats.Song {
	song, err := beats.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}
	return *song
}

func getDefaultSong() beats.Song {
	song, err := beats.Default()
	if err != nil {
		log.Fatal(err)
	}
	return *song
}
