package song

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

// Beat is all the sounds happening at a single step of the rhythm
// Step is what step of the song pattern this beat is for
type Beat struct {
	Step               int                `json:"step,omitempty"`
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
	header := fmt.Sprintf("Beat{ Step: %d", b.Step)

	inner := ""
	if b.BassDrum != bdNone {
		inner = fmt.Sprintf("%s, bd%d", inner, b.BassDrum)
	}
	if b.SnareDrum != sdNone {
		inner = fmt.Sprintf("%s, sd%d", inner, b.SnareDrum)
	}
	if b.LowTom == tOn {
		inner = fmt.Sprintf("%s, lt", inner)
	}
	if b.MidTom == tOn {
		inner = fmt.Sprintf("%s, mt", inner)
	}
	if b.HiTom == tOn {
		inner = fmt.Sprintf("%s, ht", inner)
	}
	switch b.RimshotCowbell {
	case rimshot:
		inner = fmt.Sprintf("%s, rim", inner)
	case cowbell:
		inner = fmt.Sprintf("%s, cow", inner)
	}
	switch b.HandClapTambourine {
	case handClap:
		inner = fmt.Sprintf("%s, hcp", inner)
	case tambourine:
		inner = fmt.Sprintf("%s, tamb", inner)
	}
	switch b.HiHat {
	case open:
		inner = fmt.Sprintf("%s, hho", inner)
	case closed:
		inner = fmt.Sprintf("%s, hhc", inner)
	}
	switch b.Cymbal {
	case crash:
		inner = fmt.Sprintf("%s, cyc", inner)
	case ride:
		inner = fmt.Sprintf("%s, cyr", inner)
	}
	if b.Accent == acOn {
		inner = fmt.Sprintf("%s, ac", inner)
	}

	return fmt.Sprintf("%s%s }", header, inner)
}

// ByStep implements sort.Interface for []Beat based on the Step field.
type ByStep []Beat

func (s ByStep) Len() int           { return len(s) }
func (s ByStep) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByStep) Less(i, j int) bool { return s[i].Step < s[j].Step }
