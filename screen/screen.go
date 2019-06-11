package screen

import (
	"fmt"
)

type Screen interface {
	Write(pos int, line string) bool
	Clear() bool
	Prompt()
}

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
	csi        string = "\x1B["
	gotoLine   string = csi + "%dH"
	clearLine  string = csi + "2K"
	clrScreeen string = csi + "2J"
)

var Instance = screenImpl{0, make(chan writeMsg)}

func (s *screenImpl) Write(pos int, line string) {
	s.ch <- writeMsg{pos, line, s.promptLine}
}

func (s *screenImpl) Clear() {
	s.ch <- writeMsg{0, clrScreeen, 0}
}

func (s *screenImpl) Prompt(line int) {
	s.promptLine = line
	s.ch <- writeMsg{0, fmt.Sprintf(gotoLine+clearLine+"> ", line), 0}
}

func init() {
	go func() {
		for {
			msg := <-Instance.ch
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
