package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cody-s-lee/beats/song"

	"github.com/benbjohnson/clock"
)

func main() {
	args := os.Args[1:]

	var song song.Song
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
	out := make(chan string)

	go song.Play(clock, out)

	for l := range out {
		fmt.Printf("%s\n", l)
	}
}

func getSong(reader io.Reader) song.Song {
	song, err := song.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}
	return *song
}

func getDefaultSong() song.Song {
	song, err := song.Default()
	if err != nil {
		log.Fatal(err)
	}
	return *song
}
