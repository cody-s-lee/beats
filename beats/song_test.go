package beats_test

import (
	"log"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/cody-s-lee/beats/beats"
	"github.com/google/go-cmp/cmp"
)

// TestPlay verifies that a song plays by using the default song and verifying
// it has a name, tempo and outputs appropriately when the clock advances
func TestPlay(t *testing.T) {
	song, err := beats.Default()
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
	out := make(chan beats.Step)

	go func() {
		song.Play(clock, out)
	}()

	// Advance to just before first tick, we are prepared for first tick
	advance(clock, song.TickDuration()/2)

	tick := 0
	// For each beat in the song
	for i := 0; i < len(song.Beats); i++ {
		// Advance until ready for that beat

		for ; tick < song.Beats[i].Tick; tick++ {
			// Verify that nothing's waiting on the output channel
			step, ok := read(out)
			if !ok {
				t.Fatal("Output channel should be open")
			} else if !cmp.Equal(step.Beat, beats.Beat{}) {
				t.Fatalf("Output channel should be empty but got %s\n", step.Beat)
			}
			advance(clock, song.TickDuration())
		}

		// Verify there's a beat on the output channel
		step, ok := read(out)
		if !ok {
			t.Fatal("Output channel should be open")
		} else if cmp.Equal(step.Beat, beats.Beat{}) {
			t.Fatalf("Output channel should be non-empty but got %d, %s\n", step.Tick, step.Beat)
		}
		advance(clock, song.TickDuration())
		tick++
	}

	// Verify that the output channel has closed
	step, ok := read(out)
	if ok {
		t.Fatalf("Output channel should be empty but got %s\n", step.Beat)
	}
}

// TestParseNonJson verifies we fail to parse a non-json file
func TestParseNonJson(t *testing.T) {
	reader, err := os.Open("testdata/non.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = beats.Parse(reader)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

// TestParseDuplicateTick verifies we fail to parse when we have duplicate tick value
func TestParseDuplicateTick(t *testing.T) {
	reader, err := os.Open("testdata/duplicate-tick.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = beats.Parse(reader)
	if err == nil {
		t.Fatal("Expected an error")
	}
}

//TestParseNegativeTick verifies we fail to parse when we have a negative tick value
func TestParseNegativeTick(t *testing.T) {
	reader, err := os.Open("testdata/negative-tick.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = beats.Parse(reader)
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

	_, err = beats.Parse(reader)
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

	_, err = beats.Parse(reader)
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

	_, err = beats.Parse(reader)
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

	_, err = beats.Parse(reader)
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

	_, err = beats.Parse(reader)
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

	song, err := beats.Parse(reader)
	if err != nil {
		t.Fatal(err)
	}

	clock := clock.NewMock()
	out := make(chan beats.Step)

	go func() {
		song.Play(clock, out)
	}()

	for {
		advance(clock, song.TickDuration())
		select {
		case step, ok := <-out:
			t.Log(step)
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

	song, err := beats.Parse(reader)
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
func read(out chan beats.Step) (beats.Step, bool) {
	select {
	case step, ok := <-out:
		return step, ok
	default:
		return beats.Step{}, true
	}
}
