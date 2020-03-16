package beats

import "fmt"

// Bass represents the state of the bass drum: none, drum 1 or drum 2
type Bass uint8

const (
    bdNone = Bass(iota)
    bdOne
    bdTwo
)

// Snare represents the state of the snare drum: none, drum 1 or drum 2
type Snare uint8

const (
    sdNone = Snare(iota)
    sdOne
    sdTwo
)

// RimshotCowbell represents the rimshot and cowbell sounds: none, rimshot or cowbell
type RimshotCowbell uint8

const (
    rcNone = RimshotCowbell(iota)
    rimshot
    cowbell
)

// HandClapTambourine represents the hand clap and tambourine sounds: none, hand clap or tambourine
type HandClapTambourine uint8

const (
    htNone = HandClapTambourine(iota)
    handClap
    tambourine
)

// HiHat represents the hi-hat: none, closed or open
type HiHat uint8

const (
    hhNone = HiHat(iota)
    closed
    open
)

// Cymbal represents the cymbal: none, crash or ride
type Cymbal uint8

const (
    cyNone = Cymbal(iota)
    crash
    ride
)

// Tom represents the toms: none or on
type Tom uint8

const (
    tNone = Tom(iota)
    tOn
)

// Accent represents the accent: none or on
type Accent uint8

const (
    acNone = Accent(iota)
    acOn
)

// Beat is all the sounds happening at a single tick of the rhythm
// Tick is what tick of the song pattern this beat is for
type Beat struct {
    Tick               int                `json:"tick,omitempty"`
    BassDrum           Bass               `json:"bd,omitempty"`
    SnareDrum          Snare              `json:"sd,omitempty"`
    LowTom             Tom                `json:"lt,omitempty"`
    MidTom             Tom                `json:"mt,omitempty"`
    HiTom              Tom                `json:"ht,omitempty"`
    RimshotCowbell     RimshotCowbell     `json:"rc,omitempty"`
    HandClapTambourine HandClapTambourine `json:"hc,omitempty"`
    HiHat              HiHat              `json:"hh,omitempty"`
    Cymbal             Cymbal             `json:"cy,omitempty"`
    Accent             Accent             `json:"ac,omitempty"`
}

func (b Beat) String() string {
    s := ""
    if b.BassDrum != bdNone {
        s = fmt.Sprintf("%s+bass_%d", s, b.BassDrum)
    }
    if b.SnareDrum != sdNone {
        s = fmt.Sprintf("%s+snare_%d", s, b.SnareDrum)
    }
    if b.LowTom == tOn {
        s = fmt.Sprintf("%s+low_tom", s)
    }
    if b.MidTom == tOn {
        s = fmt.Sprintf("%s+mid_tom", s)
    }
    if b.HiTom == tOn {
        s = fmt.Sprintf("%s+hi_tom", s)
    }
    switch b.RimshotCowbell {
    case rimshot:
        s = fmt.Sprintf("%s+rim", s)
    case cowbell:
        s = fmt.Sprintf("%s+cow", s)
    }
    switch b.HandClapTambourine {
    case handClap:
        s = fmt.Sprintf("%s+hcp", s)
    case tambourine:
        s = fmt.Sprintf("%s+tamb", s)
    }
    switch b.HiHat {
    case open:
        s = fmt.Sprintf("%s+hh_open", s)
    case closed:
        s = fmt.Sprintf("%s+hh_closed", s)
    }
    switch b.Cymbal {
    case crash:
        s = fmt.Sprintf("%s+cy_crash", s)
    case ride:
        s = fmt.Sprintf("%s+cy_ride", s)
    }
    if b.Accent == acOn {
        s = fmt.Sprintf("%s+acc", s)
    }

    if len(s) > 0 {
        s = s[1:]
    }

    return s
}

// ByTick implements sort.Interface for []Beat based on the Tick field.
type ByTick []Beat

func (s ByTick) Len() int           { return len(s) }
func (s ByTick) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByTick) Less(i, j int) bool { return s[i].Tick < s[j].Tick }
