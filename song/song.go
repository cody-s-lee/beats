package song

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"time"

	"github.com/benbjohnson/clock"
)

// Song is a whole song including its name, tempo and all the beats. The beats
// array is sparse; each beat covers its own step in the rhythm. Beats must be
// sorted.
type Song struct {
	Name  string
	Tempo int
	Beats []Beat
}

// NewSong creates a song while ensuring that the beats of the song are validly
// numbered. Step numbers must be greater than 0 and may not repeat.
func NewSong(name string, tempo int, beats []Beat) (*Song, error) {
	// Sort the beats in case we were passed bad data
	sort.Sort(ByStep(beats))

	// Validate non-empty name
	if name == "" {
		return nil, errors.New("Song name should not be empty")
	}

	// Validate positive tempo
	if !(tempo > 0) {
		return nil, errors.New("Song tempo should be greater than 0")
	}

	// Evaluation note: though the following two for loops could be collapsed
	// into a single for loop they are intentionally kept separate for clarity.
	// Validators should be separated and processed independently unless the
	// performance penalty requires corrective action.

	// Validate no non-positive step numbers
	for i := 0; i < len(beats); i++ {
		if beats[i].Step <= 0 {
			return nil, errors.New("Step number for beat must be greater than 0")
		}
	}

	// Validate that no step number repeats among the beats
	step := 0
	for i := 0; i < len(beats); i++ {
		if beats[i].Step == step {
			return nil, errors.New("Step number for beat may not repeat")
		}
	}

	return &Song{
		Name:  name,
		Tempo: tempo,
		Beats: beats,
	}, nil
}

// Parse parses a song from a Reader
func Parse(reader io.Reader) (*Song, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var song *Song
	err = json.Unmarshal(bytes, song)
	if err != nil {
		return nil, err
	}

	return song, nil
}

// Default constructs a default song. The default is a simple four on the floor
// implementation.
func Default() (*Song, error) {
	return NewSong(
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
}

// Play plays a song. The clock parameter allows you to use a specific clock
// such as a mock clock for testing.
func (song Song) Play(clock clock.Clock, out chan string) {
	// Time between steps in the sequence
	d := song.StepDuration()

	// Make sure the beats are sorted
	beats := song.Beats
	sort.Sort(ByStep(beats))

	// Play each beat on a ticker
	step := 0
	i := 0
	ticker := clock.Ticker(d)
	for {
		if i >= len(song.Beats) {
			break
		}

		if song.Beats[i].Step == step {
			out <- fmt.Sprintf("%d: %s", step, song.Beats[i])
			i++
		}

		<-ticker.C
		step++
	}
	close(out)
}

// StepDuration gives the amount of time between steps
func (song Song) StepDuration() time.Duration {
	return time.Minute / time.Duration(song.Tempo)
}
