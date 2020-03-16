package beats

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "sort"
    "strconv"
    "time"

    "github.com/benbjohnson/clock"
    "github.com/nsf/termbox-go"
)

type state struct {
    input      bool
    firstTick  int
    activeTick int
    cursor     bool
    song       Song
    field      field
}

// Create runs the song creation app
func Create(song Song) {
    state := state{
        input:      false,
        firstTick:  1,
        activeTick: 1,
        cursor:     false,
        song:       song,
        field:      nameField,
    }
    clock := clock.New()
    ticker := clock.Ticker(500 * time.Millisecond)
    events := make(chan termbox.Event)
    go func() {
        for {
            events <- termbox.PollEvent()
        }
    }()

    err := termbox.Init()
    if err != nil {
        fmt.Println("Failure initiating termbox")
        os.Exit(1)
    }
    defer termbox.Close()

    termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
    termbox.SetOutputMode(termbox.Output256)
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    draw(state)
    termbox.Flush()

loop:
    for {
        select {
        case <-ticker.C:
            state.cursor = !state.cursor
            update(state)
        case ev := <-events:
            switch ev.Type {
            case termbox.EventKey:
                if ev.Key == termbox.KeyCtrlS {
                    state.song.save()
                }
                if ev.Key == termbox.KeyCtrlQ {
                    break loop
                }

                dispatch(&state, &ev)
                update(state)
                termbox.Flush()
            case termbox.EventResize, termbox.EventMouse:
                update(state)
            case termbox.EventError:
                fmt.Println("Failure during termbox loop")
                os.Exit(1)
            }
        }
    }
}

func (song Song) save() {
    beats := song.Beats
    sort.Sort(ByTick(beats))

    biggestTick := 0
    btIndex := 0
    for i, b := range beats {
        if b.Tick > biggestTick {
            empty := Beat{Tick: b.Tick}
            if empty != b {
                biggestTick = b.Tick
                btIndex = i
            }
        }
    }

    beats = beats[:btIndex+1]
    song.Beats = beats

    bytes, err := json.Marshal(song)
    if err != nil {
        panic(err)
    }

    err = ioutil.WriteFile(fmt.Sprintf("%s.json", song.Name), bytes, 0644)
    if err != nil {
        panic(err)
    }
}

func update(state state) {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
    draw(state)
    termbox.Flush()
}

type field int

const (
    nameField field = iota
    tempoField
    cymbalField
    hiHatField
    hcpTambField
    rimCowField
    hiTomField
    midTomField
    lowTomField
    snareDrumField
    bassDrumField
    accentField
)

var insts = []field{
    cymbalField,
    hiHatField,
    hcpTambField,
    rimCowField,
    hiTomField,
    midTomField,
    lowTomField,
    snareDrumField,
    bassDrumField,
    accentField,
}

func (song *Song) update(beatUp *Beat) {
    for i, beatOld := range song.Beats {
        if beatUp.Tick == beatOld.Tick {
            song.Beats[i] = *beatUp
            return
        }
    }
    beats := append(song.Beats, *beatUp)
    sort.Sort(ByTick(beats))
    song.Beats = beats
}

func (song Song) on(tick int) *Beat {
    for _, beat := range song.Beats {
        if beat.Tick == tick {
            return &beat
        }
    }
    return nil
}

func draw(state state) {
    termbox.SetCell(0, 0, borderTopLeft, termbox.ColorWhite, termbox.ColorBlack)
    termbox.SetCell(79, 0, borderTopRight, termbox.ColorWhite, termbox.ColorBlack)
    termbox.SetCell(0, 23, borderBottomLeft, termbox.ColorWhite, termbox.ColorBlack)
    termbox.SetCell(79, 23, borderBottomRight, termbox.ColorWhite, termbox.ColorBlack)

    for x := 1; x < 79; x++ {
        termbox.SetCell(x, 0, borderVertical, termbox.ColorWhite, termbox.ColorBlack)
        termbox.SetCell(x, 3, borderVertical, termbox.ColorWhite, termbox.ColorBlack)
        termbox.SetCell(x, 23, borderVertical, termbox.ColorWhite, termbox.ColorBlack)
    }
    for y := 1; y < 23; y++ {
        termbox.SetCell(0, y, borderHorizontal, termbox.ColorWhite, termbox.ColorBlack)
        termbox.SetCell(79, y, borderHorizontal, termbox.ColorWhite, termbox.ColorBlack)
    }
    for y := 4; y < 23; y++ {
        termbox.SetCell(12, y, borderHorizontal, termbox.ColorWhite, termbox.ColorBlack)
    }

    for y := 4; y < 23; y = y + 2 {
        for x := 13; x < 79; x++ {
            termbox.SetCell(x, y, borderVertical, termbox.ColorBlack|termbox.AttrBold, termbox.ColorBlack)
        }
    }

    for y := 4; y < 23; y++ {
        for x := 15; x < 79; x = x + 4 {
            var ch rune
            if y%2 == 0 {
                ch = midDot
            } else {
                ch = borderHorizontal
            }

            fg := termbox.ColorBlack | termbox.AttrBold
            bg := termbox.ColorBlack

            if y%2 == 0 {
                tick := state.firstTick + ((x - 15) / 4)
                if b := state.song.on(tick); b != nil {
                    var field field
                    for f, s := range fm {
                        if s.y == y {
                            field = f
                            break
                        }
                    }

                    switch field {
                    case cymbalField:
                        switch b.Cymbal {
                        case crash:
                            ch = 'c'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        case ride:
                            ch = 'r'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        }
                    case hiHatField:
                        switch b.HiHat {
                        case open:
                            ch = 'o'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        case closed:
                            ch = 'c'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        }
                    case hcpTambField:
                        switch b.HandClapTambourine {
                        case handClap:
                            ch = 'h'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        case tambourine:
                            ch = 't'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        }
                    case rimCowField:
                        switch b.RimshotCowbell {
                        case rimshot:
                            ch = 'r'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        case cowbell:
                            ch = 'c'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        }
                    case hiTomField:
                        switch b.HiTom {
                        case tOn:
                            ch = whiteSquare
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        }
                    case midTomField:
                        switch b.MidTom {
                        case tOn:
                            ch = whiteSquare
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        }
                    case lowTomField:
                        switch b.LowTom {
                        case tOn:
                            ch = whiteSquare
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        }
                    case snareDrumField:
                        switch b.SnareDrum {
                        case sdOne:
                            ch = '1'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        case sdTwo:
                            ch = '2'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        }
                    case bassDrumField:
                        switch b.BassDrum {
                        case bdOne:
                            ch = '1'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        case bdTwo:
                            ch = '2'
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        }
                    case accentField:
                        switch b.Accent {
                        case acOn:
                            ch = whiteSquare
                            fg = termbox.ColorBlack
                            bg = termbox.ColorWhite
                        }
                    }
                }

                if fm[state.field].y == y && tick == state.activeTick {
                    fg = termbox.ColorWhite
                    bg = termbox.ColorGreen

                    if state.input {
                        bg = termbox.ColorRed
                    }
                    if state.cursor {
                        fg = fg &^ termbox.AttrBold
                    }
                }
            }

            termbox.SetCell(x, y, ch, fg, bg)
        }
    }

    printfTb(1, 1, termbox.ColorWhite, termbox.ColorBlack, "Name:")
    {
        fg := termbox.ColorWhite
        bg := termbox.ColorBlack
        if state.field == nameField {
            fg = termbox.ColorWhite
            bg = termbox.ColorGreen
            if state.input {
                bg = termbox.ColorRed
            }
            if state.cursor {
                fg = fg | termbox.AttrBold
            }
        }
        printfTb(7, 1, fg, bg, "%-60s", state.song.Name)
    }

    printfTb(68, 1, termbox.ColorWhite, termbox.ColorBlack, "Tempo:")
    {
        fg := termbox.ColorWhite
        bg := termbox.ColorBlack
        if state.field == tempoField {
            fg = termbox.ColorWhite
            bg = termbox.ColorGreen
            if state.input {
                bg = termbox.ColorRed
            }
            if state.cursor {
                fg = fg | termbox.AttrBold
            }
        }
        printfTb(75, 1, fg, bg, "%3d", state.song.Tempo)
    }

    printfTb(1, 2, termbox.ColorWhite, termbox.ColorBlack, "Step")

    for x, t := 13, state.firstTick; x < 79-4; x, t = x+4, t+1 {
        printfTb(x, 2, termbox.ColorWhite, termbox.ColorBlack, fmt.Sprintf("%3d", t))
    }

    for _, i := range insts {
        fg := termbox.ColorWhite
        bg := termbox.ColorBlack
        if state.field == i {
            if state.cursor {
                fg = fg | termbox.AttrBold
            }
            fg = fg | termbox.AttrReverse
            bg = bg | termbox.AttrReverse
        }
        printfTb(fm[i].x, fm[i].y, fg, bg, "%-11s", fm[i].name)
    }
}

type fs struct {
    name  string
    x     int
    y     int
    left  field
    right field
    up    field
    down  field
}

var fm = map[field]fs{
    nameField:      fs{"Name:", 1, 1, nameField, tempoField, nameField, cymbalField},
    tempoField:     fs{"Tempo:", 68, 1, nameField, tempoField, tempoField, cymbalField},
    cymbalField:    fs{"CYmbal", 1, 4, cymbalField, cymbalField, nameField, hiHatField},
    hiHatField:     fs{"HiHat", 1, 6, hiHatField, hiHatField, cymbalField, hcpTambField},
    hcpTambField:   fs{"HCP/TAMB", 1, 8, hcpTambField, hcpTambField, hiHatField, rimCowField},
    rimCowField:    fs{"RIM/COWbell", 1, 10, rimCowField, rimCowField, hcpTambField, hiTomField},
    hiTomField:     fs{"Hi Tom", 1, 12, hiTomField, hiTomField, rimCowField, midTomField},
    midTomField:    fs{"Mid Tom", 1, 14, midTomField, midTomField, hiTomField, lowTomField},
    lowTomField:    fs{"Low Tom", 1, 16, lowTomField, lowTomField, midTomField, snareDrumField},
    snareDrumField: fs{"Snare Drum", 1, 18, snareDrumField, snareDrumField, lowTomField, bassDrumField},
    bassDrumField:  fs{"Bass Drum", 1, 20, bassDrumField, bassDrumField, snareDrumField, accentField},
    accentField:    fs{"ACcent", 1, 22, accentField, accentField, bassDrumField, accentField},
}

func (beat *Beat) normalize() {
    if beat.Cymbal > ride {
        beat.Cymbal = ride
    }
    if beat.BassDrum > bdTwo {
        beat.BassDrum = bdTwo
    }
    if beat.SnareDrum > sdTwo {
        beat.SnareDrum = sdTwo
    }
    if beat.RimshotCowbell > cowbell {
        beat.RimshotCowbell = cowbell
    }
    if beat.HandClapTambourine > tambourine {
        beat.HandClapTambourine = tambourine
    }
    if beat.HiHat > open {
        beat.HiHat = open
    }
    if beat.Accent > acOn {
        beat.Accent = acOn
    }
    if beat.LowTom > tOn {
        beat.LowTom = tOn
    }
    if beat.MidTom > tOn {
        beat.MidTom = tOn
    }
    if beat.HiTom > tOn {
        beat.HiTom = tOn
    }
}

func (beat *Beat) set(field field, value int) {
    if value < 0 {
        value = 0
    }
    switch field {
    case cymbalField:
        beat.Cymbal = Cymbal(value)
    case hiHatField:
        beat.HiHat = HiHat(value)
    case hcpTambField:
        beat.HandClapTambourine = HandClapTambourine(value)
    case rimCowField:
        beat.RimshotCowbell = RimshotCowbell(value)
    case hiTomField:
        beat.HiTom = Tom(value)
    case midTomField:
        beat.MidTom = Tom(value)
    case lowTomField:
        beat.LowTom = Tom(value)
    case snareDrumField:
        beat.SnareDrum = Snare(value)
    case bassDrumField:
        beat.BassDrum = Bass(value)
    case accentField:
        beat.Accent = Accent(value)
    }

    beat.normalize()
}

func (beat Beat) value(field field) int {
    switch field {
    case cymbalField:
        return int(beat.Cymbal)
    case hiHatField:
        return int(beat.HiHat)
    case hcpTambField:
        return int(beat.HandClapTambourine)
    case rimCowField:
        return int(beat.RimshotCowbell)
    case hiTomField:
        return int(beat.HiTom)
    case midTomField:
        return int(beat.MidTom)
    case lowTomField:
        return int(beat.LowTom)
    case snareDrumField:
        return int(beat.SnareDrum)
    case bassDrumField:
        return int(beat.BassDrum)
    case accentField:
        return int(beat.Accent)
    }
    return -1
}

func dispatch(state *state, ev *termbox.Event) {
    if state.input {
        switch state.field {
        case nameField:
            if ev.Ch != 0 {
                state.song.Name = state.song.Name + string(ev.Ch)
            } else if len(state.song.Name) > 0 && ev.Key == termbox.KeyBackspace {
                state.song.Name = state.song.Name[:len(state.song.Name)-1]
            } else if ev.Key == termbox.KeySpace {
                state.song.Name = state.song.Name + " "
            }
        case tempoField:
            if ev.Ch != 0 {
                i, err := strconv.Atoi(string(ev.Ch))
                if err == nil {
                    state.song.Tempo = state.song.Tempo*10 + i
                }
            } else if ev.Key == termbox.KeyBackspace {
                state.song.Tempo = int(state.song.Tempo / 10)
            }
        default:
            beat := state.song.on(state.activeTick)
            if beat == nil {
                beat = &Beat{
                    Tick: state.activeTick,
                }
            }

            if ev.Ch == 0 {
                switch ev.Key {
                case termbox.KeyArrowLeft, termbox.KeyArrowUp:
                    beat.set(state.field, beat.value(state.field)-1)
                case termbox.KeyArrowRight, termbox.KeyArrowDown:
                    beat.set(state.field, beat.value(state.field)+1)
                }
            }

            termbox.Flush()
            state.song.update(beat)
        }

        if ev.Key == termbox.KeyEnter {
            state.input = !state.input
            if state.field == tempoField {
                if state.song.Tempo <= 0 {
                    state.song.Tempo = 1
                }
            }
        }
    } else {
        switch ev.Key {
        case termbox.KeyArrowLeft:
            state.field = fm[state.field].left
            if state.field != nameField && state.field != tempoField {
                state.activeTick--
                if state.activeTick < 1 {
                    state.activeTick = 1
                }
            }
        case termbox.KeyArrowRight:
            state.field = fm[state.field].right
            if state.field != nameField && state.field != tempoField {
                state.activeTick++
            }
        case termbox.KeyArrowUp:
            state.field = fm[state.field].up
        case termbox.KeyArrowDown:
            state.field = fm[state.field].down
        case termbox.KeyEnter:
            state.input = !state.input
        }

        state.firstTick = state.activeTick - 7
        if state.firstTick < 1 {
            state.firstTick = 1
        }
    }
}

func printTb(x, y int, fg, bg termbox.Attribute, msg string) {
    for _, c := range msg {
        termbox.SetCell(x, y, c, fg, bg)
        x++
    }
}

func printfTb(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) {
    s := fmt.Sprintf(format, args...)
    printTb(x, y, fg, bg, s)
}

const borderTopLeft rune = 0x250C
const borderTopRight rune = 0x2510
const borderBottomLeft rune = 0x2514
const borderBottomRight rune = 0x2518
const borderVertical rune = 0x2500
const borderHorizontal rune = 0x2502
const borderHorizontalLeftBar rune = 0x251C
const borderHorizontalRight rune = 0x2524
const boxShadow rune = 0x2588
const midDot rune = 0x00B7
const blackSquare rune = 0x25A0
const whiteSquare rune = 0x25A1
