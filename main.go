package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/benbjohnson/clock"
	"github.com/cody-s-lee/beats/beats"
)

func main() {
	args := os.Args[1:]

	// No arguments given: play default song and quit gracefully
	if len(args) == 0 {
		song := getDefaultSong()
		play(song)
		os.Exit(0)
	}

	switch args[0] {
	case "play":
		// Not enough args for play, show help and quit
		if len(args) < 2 {
			showHelp()
			os.Exit(1)
		}

		// Grab file
		fn := args[1]
		reader, err := os.Open(fn)
		if err != nil {
			fmt.Printf("Could not open file %s\n", fn)
			showHelp()
			os.Exit(1)
		}

		// Play file
		song := getSong(reader)
		play(song)
		os.Exit(0)

	case "create":
		song := beats.Song{Tempo: 100}

		if len(args) > 1 {
			// Grab file
			fn := args[1]
			reader, err := os.Open(fn)
			if err != nil {
				fmt.Printf("Could not open file %s\n", fn)
				showHelp()
				os.Exit(1)
			}

			// Load song from file
			song = getSong(reader)
		}

		create(song)
		os.Exit(0)

	case "help", "-h", "--help":
		showHelp()
		os.Exit(0)
	}

	// No valid command given
	showHelp()
	os.Exit(1)
}

func showHelp() {
	fmt.Printf(
		`usage: %s <command> [<args>]

command is one of:
    play <filename>        Play a song
    create [filename]      Create a song


If no command is given the default song (four on the floor) is played.

Play Mode:

play takes a filename as an argument to load a song from a file. The song file should be a json representation of the song in the following format:

{
    "name": "song name",
    "tempo": 100,
    "beats": [ <beat>... ]
}

- song name is required and must be non-empty
- tempo is required and must be positive, denoted in beats per minute (bpm)
- beats is an array of beat objects of the following format:

{
    "tick": 1,
    "bd": 0,
    "sd": 0,
    "lt": 0,
    "mt": 0,
    "ht": 0,
    "rc": 0,
    "hc": 0,
    "hh": 0,
    "cy": 0,
    "ac": 0
}

- tick is required and must be positive, denotes which tick the note occurs on, starting with 1
- The rest of the fields are individual instruments using integer assignments. Fields can be omitted if the note is inactive on this beat. Instruments are inactive for a value of 0.
-- bd: Bass drum           - off (0), drum 1 (1), drum 2 (2)
-- sd: Snare drum          - off (0), drum 1 (1), drum 2 (2)
-- lt: Low tom             - off (0), active (1)
-- mt: Mid tom             - off (0), active (1)
-- ht: High tom            - off (0), active (1)
-- rc: Rimshot/Cowbell     - off (0), rimshot (1), cowbell (2)
-- hc: Handclap/Tambourine - off (0), rimshot (1), cowbell (2)
-- hh: Hi-Hat              - off (0), closed (1), open (2)
-- cy: Cymbal              - off (0), crash (1), ride (2)
-- ac: Accent              - off (0), active (1)

Create Mode:

create has a term-based ui for song creation. Optionally a filename of a song can be used to load in a song to work on.

Commands:
    ctrl-s to save to <name>.json
    ctrl-q to quit

    enter to enter or leave input mode for highlighted cell
    arrow keys modify the current cell when in input mode
    arrow keys move around the board when not in input mode
`, os.Args[0])
}

func create(song beats.Song) {
	beats.Create(song)
}

func play(song beats.Song) {
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
