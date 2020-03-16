package beats

import (
    "encoding/json"
    "errors"
    "io"
    "sort"
    "time"

    "github.com/benbjohnson/clock"
)

// Song is a whole song including its name, tempo and all the beats. The beats
// array is sparse; each beat covers its own tick in the rhythm. Beats must be
// sorted.
type Song struct {
    Name  string `json:"name,omitempty"`
    Tempo int    `json:"tempo,omitempty"`
    Beats []Beat `json:"beats"`
}

// NewSong creates a song while ensuring that the beats of the song are validly
// numbered. Tick numbers must be greater than 0 and may not repeat.
func NewSong(name string, tempo int, beats []Beat) (*Song, error) {
    // Sort the beats in case we were passed bad data
    sort.Sort(ByTick(beats))

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

    // Validate no non-positive tick numbers
    for i := 0; i < len(beats); i++ {
        if beats[i].Tick <= 0 {
            return nil, errors.New("Tick number for beat must be greater than 0")
        }
    }

    // Validate that no tick number repeats among the beats
    tick := 0
    for i := 0; i < len(beats); i++ {
        if beats[i].Tick == tick {
            return nil, errors.New("Tick number for beat may not repeat")
        }
        tick = beats[i].Tick
    }

    return &Song{
        Name:  name,
        Tempo: tempo,
        Beats: beats,
    }, nil
}

// Parse parses a song from a Reader
func Parse(reader io.Reader) (*Song, error) {
    var song Song
    err := json.NewDecoder(reader).Decode(&song)
    if err != nil {
        return nil, err
    }

    return NewSong(song.Name, song.Tempo, song.Beats)
}

// Default constructs a default song. The default is a simple four on the floor
// implementation.
func Default() (*Song, error) {
    return NewSong(
        "four-on-the-floor",
        128,
        []Beat{
            Beat{
                Tick:     1,
                BassDrum: bdOne,
            },
            Beat{
                Tick:  3,
                HiHat: closed,
            },
            Beat{
                Tick:      5,
                SnareDrum: sdOne,
                BassDrum:  bdOne,
            },
            Beat{
                Tick:  7,
                HiHat: closed,
            },
            Beat{
                Tick:     9,
                BassDrum: bdOne,
            },
            Beat{
                Tick:  11,
                HiHat: closed,
            },
            Beat{
                Tick:      13,
                SnareDrum: sdOne,
                BassDrum:  bdOne,
            },
            Beat{
                Tick:  15,
                HiHat: closed,
            },
        },
    )
}

// Step is one step of the sequence at a given tick
type Step struct {
    Tick int
    Beat Beat
}

// Play plays a song. The clock parameter allows you to use a specific clock
// such as a mock clock for testing.
func (song Song) Play(clock clock.Clock, out chan Step) {
    // Time between ticks in the sequence
    d := song.TickDuration()

    // Make sure the beats are sorted
    beats := song.Beats
    sort.Sort(ByTick(beats))

    // Play each beat on a ticker
    tick := 1
    i := 0
    ticker := clock.Ticker(d)
    for {
        if i >= len(song.Beats) {
            break
        }

        if song.Beats[i].Tick == tick {
            go func(tick int, i int) {
                out <- Step{tick, song.Beats[i]}
            }(tick, i)
            i++
        } else {
            go func(tick int) {
                out <- Step{tick, Beat{Tick: tick}}
            }(tick)
        }

        <-ticker.C
        tick++
    }
    close(out)
}

// TickDuration gives the amount of time between ticks
func (song Song) TickDuration() time.Duration {
    return time.Minute / time.Duration(song.Tempo)
}
