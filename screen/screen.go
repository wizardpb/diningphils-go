package screen

import (
	"fmt"
)

var silent = false

type writeMsg struct {
	pos        int
	line       string
	promptLine int
}

type screenImpl struct {
	promptLine int
	ch         chan writeMsg
}

const (
	csi       string = "\x1B["
	gotoLine  string = csi + "%dH"
	clearLine string = csi + "2K"
	clrScreen string = csi + "2J"
)

var screen = screenImpl{0, make(chan writeMsg)}

func Silent(on bool) {
	silent = on
}

func Write(pos int, line string) {
	if !silent {
		screen.ch <- writeMsg{pos, line, screen.promptLine}
	}
}

func Clear() {
	screen.ch <- writeMsg{0, fmt.Sprintf(clrScreen+gotoLine, 1), 0}
}

func Prompt(line int) {
	screen.promptLine = line // Simply because save and restore cursor don't work
	screen.ch <- writeMsg{0, fmt.Sprintf(gotoLine+clearLine+"> ", line), 0}
}

func init() {
	go func() {
		for {
			msg := <-screen.ch
			if msg.pos > 0 {
				fmt.Printf(gotoLine+clearLine+msg.line, msg.pos)
				if msg.promptLine > 0 {
					fmt.Printf(csi+"%d;%dH", msg.promptLine, 3)
				}
			} else {
				fmt.Print(msg.line)
			}
		}
	}()
}
