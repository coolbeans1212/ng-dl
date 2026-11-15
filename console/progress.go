package console

import (
	"bytes"
	"github.com/mattn/go-colorable"
	"golang.org/x/term"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type progress struct {
	Id      int
	Action  string
	Percent int
}

var progresses = make(map[int]progress)
var progressMu sync.Mutex

var progressesBefore int

func maxInt(a int, b int) int { // grrr max() was only built-in after Go 1.21... but i don't call it max() so that it doesn't conflict if you try to build it on a later version :)
	if a > b {
		return a
	} else {
		return b
	}
}

func minInt(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

// Source - https://stackoverflow.com/a/61164288
// Posted by Robert J., modified by community. See post 'Timeline' for change history
// Retrieved 2026-11-15, License - CC BY-SA 4.0
func RandBool() bool {
	return rand.Float32() > 0.5
}

func GetTermWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return 80
	}
	return width
}

func RefreshProgress() {
	const BACKGROUND = "44" // cyan
	const FOREGROUND = "33" // yellow
	var buf bytes.Buffer
	out := colorable.NewColorableStdout()
	buf.WriteString("\033[s") // save cursor pos

	// clear the previous progress bars
	buf.WriteString("\033[1;1H") // set cursor to (1, 1)
	for i := 0; i < progressesBefore*3+3; i++ {
		buf.WriteString("\033[2K")
		if i < progressesBefore*3+2 {
			buf.WriteString("\033[1B")
		}
	}

	// draw the next progress bars
	buf.WriteString("\033[1;1H") // set cursor pos to (1, 1)
	var lines int
	if len(progresses) > 3 {
		lines = 3*3 + 2
	} else {
		lines = len(progresses)*3 + 2
	}
	if len(progresses) > 0 {
		ids := make([]int, 0, len(progresses))

		for id := range progresses {
			ids = append(ids, id)
		}

		sort.Ints(ids)

		limit := minInt(len(ids), 4)
		buf.WriteString("\033[2K\n") // blank row
		for i := 0; i < lines; i++ {
			buf.WriteString("\033[" + FOREGROUND + ";" + BACKGROUND + "m\033[K\033[0m\n") // colour in rows in with cyan
		}
		buf.WriteString("\033[3;1H") // set cursor pos to (3, 1)
		buf.WriteString("\033[" + BACKGROUND + "m\033[" + FOREGROUND + "m")
		for _, i := range ids[:limit] {
			buf.WriteString("\033[" + BACKGROUND + "m\033[" + FOREGROUND + "m " + progresses[i].Action + "\n")
			if progresses[i].Percent != -1 {
				buf.WriteString("\033[" + BACKGROUND + "m\033[" + FOREGROUND + "m    " + strconv.Itoa(progresses[i].Percent) + "% complete\n")
			} else {
				buf.WriteString("\033[" + BACKGROUND + "m\033[" + FOREGROUND + "m    Processing...\n")
			}
			buf.WriteString("\033[" + BACKGROUND + "m\033[" + FOREGROUND + "m    [")
			inner := GetTermWidth() - 8
			if progresses[i].Percent != -1 {
				if inner < 0 {
					inner = 0
				}
				filled := inner * progresses[i].Percent / 100
				if filled < 0 {
					filled = 0
				}
				buf.WriteString(strings.Repeat("o", filled))
				buf.WriteString(strings.Repeat(" ", inner-filled))
			} else {
				for i := 0; i < inner; i++ {
					if RandBool() {
						buf.WriteString("?")
					} else {
						buf.WriteString(" ")
					}
				}
			}
			buf.WriteString("]\n")
		}
		if len(progresses) > 3 {
			buf.WriteString(strconv.Itoa(len(progresses)-4) + " activities not shown...")
		}
		progressesBefore = len(progresses)
	} else {
		progressesBefore = 0
	}
	buf.WriteString("\033[0m") // reset text formatting
	buf.WriteString("\033[u")  // restore cursor pos
	out.Write(buf.Bytes())
}

func WriteProgress(id int, action string, percent int) {
	progressMu.Lock()
	defer progressMu.Unlock()
	progresses[id] = progress{Id: id, Action: action, Percent: percent}
	RefreshProgress()
}

func ClearProgress(id int) {
	progressMu.Lock()
	defer progressMu.Unlock()
	delete(progresses, id)
	RefreshProgress()
}
