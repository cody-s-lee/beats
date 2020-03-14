package song_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/cody-s-lee/beats/song"
)

func TestPlay(t *testing.T) {
	song, err := song.Default()
	if err != nil {
		t.Error(err)
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
		t.Errorf("Output channel should be empty but got %s\n", l)
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
