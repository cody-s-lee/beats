package song_test

import (
	"log"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/cody-s-lee/beats/song"
)

// TestPlay verifies that a song plays by using the default song and verifying
// it has a name, tempo and outputs appropriately when the clock advances
func TestPlay(t *testing.T) {
	song, err := song.Default()
	if err != nil {
		t.Fatal(err)
	}

	if song.Name == "" {
		t.Fatal("Song name should be non-empty")
	}

	if !(song.Tempo > 0) {
		t.Fatal("Song tempo should be greater than 0")
	}

	clock := clock.NewMock()
	out := make(chan string)

	go func() {
		song.Play(clock, out)
	}()

	// Advance to just before first step, we are prepared for first step
	advance(clock, song.StepDuration()/2)
	step := 1

	// For each beat in the song
	for i := 0; i < len(song.Beats); i++ {
		// Advance until ready for that beat
		for ; step < song.Beats[i].Step; step++ {
			advance(clock, song.StepDuration())
		}

		// Verify that nothing's waiting on the output channel
		l, ok := read(out)
		if ok {
			t.Fatalf("Output channel should be empty but got %s\n", l)
		}

		// Advance to the next step
		advance(clock, song.StepDuration())
		step++

		// Verify there's a beat on the output channel
		l, ok = read(out)
		if ok {
			t.Log(l)
		} else {
			t.Fatal("Output channel should be populated")
		}
	}

	// Advance past the last step
	advance(clock, song.StepDuration())
	step++

	// Verify that nothing's waiting on the output channel
	l, ok := read(out)
	if ok {
		t.Fatalf("Output channel should be empty but got %s\n", l)
	}
}

// TestParseNonJson verifies we fail to parse a non-json file
func TestParseNonJson(t *testing.T) {
	reader, err := os.Open("testdata/non.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = song.Parse(reader)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

// TestParseDuplicateStep verifies we fail to parse when we have duplicate step value
func TestParseDuplicateStep(t *testing.T) {
	reader, err := os.Open("testdata/duplicate-step.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = song.Parse(reader)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

//TestParseNegativeStep verifies we fail to parse when we have a negative step value
func TestParseNegativeStep(t *testing.T) {
	reader, err := os.Open("testdata/negative-step.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = song.Parse(reader)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

// TestParseEmptyName verifies we fail to parse when we have an empty song name
func TestParseEmptyName(t *testing.T) {
	reader, err := os.Open("testdata/empty-name.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = song.Parse(reader)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

// TestParseNegativeTempo verifies we fail to parse when we have a negative tempo
func TestParseNegativeTempo(t *testing.T) {
	reader, err := os.Open("testdata/negative-tempo.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = song.Parse(reader)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

// TestParseNoName verifies we fail to parse when we have no song name
func TestParseNoName(t *testing.T) {
	reader, err := os.Open("testdata/no-name.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = song.Parse(reader)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

// TestParseNoTempo verifies we fail to parse when the song has no tempo
func TestParseNoTempo(t *testing.T) {
	reader, err := os.Open("testdata/no-tempo.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = song.Parse(reader)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

// TestParseZeroTempo verifies we fail to parse when the song has a zero tempo
func TestParseZeroTempo(t *testing.T) {
	reader, err := os.Open("testdata/zero-tempo.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = song.Parse(reader)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

// TestParseAllNotes verifies we can read all valid notes
func TestParseAllNotes(t *testing.T) {
	reader, err := os.Open("testdata/all-notes.json")
	if err != nil {
		log.Fatal(err)
	}

	song, err := song.Parse(reader)
	if err != nil {
		t.Fatal(err)
	}

	clock := clock.NewMock()
	out := make(chan string)

	go func() {
		song.Play(clock, out)
	}()

	for {
		advance(clock, song.StepDuration())
		select {
		case l, ok := <-out:
			t.Log(l)
			if !ok {
				return
			}
		default:
		}
	}
}

// TestParse verifies a well-formatted file parses
func TestParse(t *testing.T) {
	reader, err := os.Open("testdata/cowbell.json")
	if err != nil {
		log.Fatal(err)
	}

	song, err := song.Parse(reader)
	if err != nil {
		t.Fatal(err)
	}

	if song.Name != "Fast Cowbell" {
		t.Errorf("Expected song name Fast Cowbell but got %s", song.Name)
	}

	if song.Tempo != 188 {
		t.Errorf("Expected tempo of 188 but got %d", song.Tempo)
	}

	if len(song.Beats) != 10 {
		t.Errorf("Expected to find 10 beats but got %d", len(song.Beats))
	}
}

// advance advances the given clock by the given duration and then
// cooperatively yield to give the player a chance to work.
func advance(clock *clock.Mock, duration time.Duration) {
	clock.Add(duration)
	runtime.Gosched()
}

// read reads off the output channel in a non-blocking manner
func read(out chan string) (string, bool) {
	select {
	case l, ok := <-out:
		return l, ok
	default:
		return "", false
	}
}
