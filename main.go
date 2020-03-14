package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"

	"github.com/benbjohnson/clock"
)

func main() {
	args := os.Args[1:]

	var song Song
	if len(args) > 0 {
		reader, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}

		song = getSong(reader)
	} else {
		song = getDefaultSong()
	}

	clock := clock.New()
	Play(song, clock)
}

// Play plays a song. The clock parameter allows you to use a specific clock
// such as a mock clock for testing.
func Play(song Song, clock clock.Clock) {
	fmt.Printf("Name: %s\n", song.Name)
	fmt.Printf("Tempo: %d bpm\n", song.Tempo)

	// Time between steps in the sequence
	d := time.Duration(time.Minute / time.Duration(song.Tempo))

	// done := make(chan struct{})
	// go func() {
	ticker := clock.Ticker(d)
	step := 0
	i := 0

	// Make sure the beats are sorted
	beats := song.Beats
	sort.Sort(ByStep(beats))
	for {
		if i >= len(song.Beats) {
			break
		}

		if song.Beats[i].Step == step {
			fmt.Printf("%d: %s\n", step, song.Beats[i])
			i++
		}

		<-ticker.C
		step++
	}
	// close(done)
	// }()
	// runtime.Gosched()

	// <-done
}

func getSong(reader io.Reader) Song {
	song, err := ParseSong(reader)
	if err != nil {
		log.Fatal(err)
	}
	return *song
}

func getDefaultSong() Song {
	song, err := NewSong(
		"four-on-the-floor",
		128,
		[]Beat{
			Beat{
				Step:     1,
				BassDrum: bdOne,
			},
			Beat{
				Step:  3,
				HiHat: closed,
			},
			Beat{
				Step:      5,
				SnareDrum: sdOne,
				BassDrum:  bdOne,
			},
			Beat{
				Step:  7,
				HiHat: closed,
			},
			Beat{
				Step:     9,
				BassDrum: bdOne,
			},
			Beat{
				Step:  11,
				HiHat: closed,
			},
			Beat{
				Step:      13,
				SnareDrum: sdOne,
				BassDrum:  bdOne,
			},
			Beat{
				Step:  15,
				HiHat: closed,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return *song
}
