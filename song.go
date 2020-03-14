package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"sort"
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

// ParseSong parses a song from a Reader
func ParseSong(reader io.Reader) (*Song, error) {
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
