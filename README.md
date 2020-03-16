# beats

A Roland TR-707 style drum machine. beats visualizes the output of a song in realtime. By default beats plays a simple four-on-the-floor pattern at 128 bpm.

# Usage

beats has a default mode and two subcommands, *play* and *create*. Please build via `go build` and run the executable to see the basic implementation of four-on-the-floor.

Tests run via `go test ./...` and cover the `song.go` and `beat.go` classes. `creator.go` is not covered for reasons later discussed.

> Note: All development and testing was done on Windows. Viability on OSX and Linux are unknown.

## Default

In default mode four-on-the-floor is played. This is visualized in the same way as play mode.

## Play

Play mode visualizes the song passed in as the argument to `beats play <filename>`. The file played should be a song file in json format.

Play mode first outputs the song name and tempo:

```
Name: four-on-the-floor
Tempo: 128 bpm
```

Then, at the tempo of the song, each beat is displayed, prefixed by the beat number. Empty beats are displayed to make it easier to keep time when watching.

```
1: bass_1
2:
3: hh_closed
4:
5: bass_1+snare_1
6:
7: hh_closed
8:
9: bass_1
10:
11: hh_closed
12:
13: bass_1+snare_1
14:
15: hh_closed
```

## Create

Create mode uses [nsf/termbox-go](https://github.com/nsf/termbox-go) to create an interactive user interface for song creation.

Optionally a filename can be passed to load in a song for editing using the command `beats create <filename>`.

> Note: Minimal error handling is completed here. If an invalid song file is loaded a json parsing failure is given.

### UI

The create mode UI is a curses-style UI. The song name and tempo are clearly listed at the top, then the steps in the song. A row for each instrument is shown and a cross-grid for locating notes. The UI makes use of colors to distinguish notes more clearly.

```
┌──────────────────────────────────────────────────────────────────────────────┐
│Name: Fast Cowbell                                                 Tempo: 188 │
│Step          1   2   3   4   5   6   7   8   9  10  11  12  13  14  15  16   │
│──────────────────────────────────────────────────────────────────────────────│
│CYmbal     │──·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───│
│           │  │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │
│HiHat      │──·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───│
│           │  │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │
│HCP/TAMB   │──·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───│
│           │  │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │
│RIM/COWbell│──·───·───c───·───·───·───c───·───·───·───c───·───·───·───c───·───│
│           │  │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │
│Hi Tom     │──·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───│
│           │  │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │
│Mid Tom    │──·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───│
│           │  │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │
│Low Tom    │──·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───│
│           │  │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │
│Snare Drum │──·───·───·───·───1───1───·───·───·───·───·───·───1───1───·───·───│
│           │  │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │
│Bass Drum  │──1───·───·───·───·───·───·───·───1───·───·───·───·───·───·───·───│
│           │  │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │   │
│ACcent     │──·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───·───│
└──────────────────────────────────────────────────────────────────────────────┘ 
```

The active field (name, tempo, instrument) is highlighted in green. The enter key toggles input mode and changes the highlight to red. The name and tempo fields can accept typed input. For each instrument field a specific step in the song is highlighted

> Note: No effort was made to limit name and tempo fields to fit in the space given. Future development should concern itself with those considerations.

#### Commands

* **ctrl-q** quits 
* **ctrl-s** saves the song to `<song name>.json`
* **enter** toggles input mode
* **arrow keys** traverse the UI when not in input mode
* **arrow keys** alter instrument settings in input mode

#### Input mode

When in input mode the selected field is highlighted in red instead of green.

##### Name and Tempo 

Name and tempo are typing fields. Tempo is limited to numeric input and only accepts positive integers.

##### Instruments

Instruments are changed via the arrow keys. Up and left decrements the instrument value, down and right increments.

> Note: This UI was only tested in Windows Command Prompt using default colors. The terminal window was larger than 80 columns by 24 rows. Attempting create mode on other systems would require testing and handling degradation due to window size changes would also be advised.

# Architecture

## Song and Beat Formats

The `Song` format is a simple go struct of a name, tempo and an array of `Beat`s. The `Beat` struct contains the step number and a field for each instrument type. When transfering to json format the instrument fields in `Beat` use abbreviations. This choice was to make it easier to manually construct a json file.

The `Song` struct could probably have been an unexported struct with all creations enforced through `NewSong`.

> Note: `Tick` was chosen for the name for step-related variables because it more closely meshes with the timing concepts used. "Beat" or "step" could have been used for variable names instead with little effect.

This format is not ideal for further extension. Many of the functions end up in a large switch or if block for handling each instrument type. A better representation may have been a map of instrument constant to boolean to make iteration and extension easier. For instruments with mutually exclusive values, such as rimshot/cowbell, additional logic would be needed to ensure that such rules are followed. The form used in this implementation example ensures that those rules are baked in to the product. This decision was made early in the process and did negatively affect the stringifier for `Beat` and later the `creator` implementation.

## Song Playing

See `song.go` lines 128-158. Playing is accomplished by passing a clock and output channel to the song object. The clock is externalized in order to allow for testing using fake clocks as provided by [benbjohnson/clock](https://github.com/benbjohnson/clock). In retrospect rather than pass the channel in to the play function, creating the channel within the play function and returning it, with the play operation happening in a goroutine is probably more idiomatically go style.

## Tests

Tests were only implemented for the `beats` package files `song.go` and `beat.go`. These are the major business logic files. `main.go` is not included within tests because its purpose is not to create a usable library piece but to interface with the user. Similarly, `creator.go` is untested though it does include some testable functions such as `update`, `on`, `normalize`, `set`, and `value`.

## Creator

The creator is an interactive UI for creating and editing songs. It was written partially as an exploration of [termbox](https://github.com/nsf/termbox-go). The creator is where it is most clear that the decision for how to structure the `Beat` struct is most lacking in usability. The `fs` struct and `fm` map were generated to help alleviate the issues but would have been well served by a better system for representing instrument objects.

The organization of `creator.go` could also be improved. As it stands it was made without an overall design in mind and built in an as-needed manner. Vitally, clearing up the various UI generation snippets into smaller chunks of code is the primary concern. Giving feedback on save would be a useful feature, perhaps through a modal dialog structure. Cleaning up the way that instruments are handled would also improve readibility by squashing the large number of `switch` blocks kicking around.

> Note: Much of `creator.go`'s interaction with `termbox` was cribbed from termbox's [`_demos/output.go`](https://github.com/nsf/termbox-go/blob/master/_demos/output.go) demo.

Also, continued development of the creator would allow for an in-UI player to be built with each step highlighted as the song progresses.

----

> Final Note: Assumptions abound about the operation of the TR-707 and similar drum machines. It's safe to assume errors in concepts around operation to be due to unfamiliarity by the author.