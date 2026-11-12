package console

import (
	"github.com/mattn/go-colorable"
	"golang.org/x/term"
	"os"
	"strconv"
	"strings"
)

type progress struct {
	Id      int
	Action  string
	Percent int
}

var progresses = make(map[int]progress)

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
	out := colorable.NewColorableStdout()
	if len(progresses) == 0 {
		return
	}
	out.Write([]byte("\033[s"))    // save cursor pos
	out.Write([]byte("\033[1;1H")) // set cursor pos to (1, 1)
	out.Write([]byte("\033[2K\n")) // blank row
	for i := 0; i < len(progresses)*3+2; i++ {
		out.Write([]byte("\033[" + FOREGROUND + ";" + BACKGROUND + "m\033[K\033[0m\n")) // colour in rows in with cyan
	}
	out.Write([]byte("\033[3;1H")) // set cursor pos to (3, 1)
	out.Write([]byte("\033[" + BACKGROUND + "m\033[" + FOREGROUND + "m"))
	for i := range progresses {
		out.Write([]byte("\033[" + BACKGROUND + "m\033[" + FOREGROUND + "m " + progresses[i].Action + "\n"))
		out.Write([]byte("\033[" + BACKGROUND + "m\033[" + FOREGROUND + "m    " + strconv.Itoa(progresses[i].Percent) + "% complete\n"))
		out.Write([]byte("\033[" + BACKGROUND + "m\033[" + FOREGROUND + "m    ["))
		inner := GetTermWidth() - 8
		if inner < 0 {
			inner = 0
		}
		filled := inner * progresses[i].Percent / 100
		if filled < 0 {
			filled = 0
		}
		out.Write([]byte(strings.Repeat("o", filled)))
		out.Write([]byte(strings.Repeat(" ", inner-filled)))
		out.Write([]byte("]\n"))
	}
	out.Write([]byte("\033[0m"))
	out.Write([]byte("\033[u")) // restore cursor pos
}

func WriteProgress(id int, action string, percent int) {
	progresses[id] = progress{Id: id, Action: action, Percent: percent}
	RefreshProgress()
}

func ClearProgress(id int) {
	delete(progresses, id)
}
