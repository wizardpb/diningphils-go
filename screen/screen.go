package screen

import (
	"fmt"
	"io"
	"os"
)

// Screen commands
type screenCmd interface {
	writeOn(thisScreen *screenImpl)
}

// Move the cursor
type cursorPos struct {
	row, col int
}

// Write at a given line, saving and restoring the cursor
type writeScreenLine struct {
	cursor cursorPos
	str    string
}

// Write a string at the current location, update the cursor pos
type writeStr struct {
	str string
}

// Clear the screen by writing a clear screen CSI and resetingthe cursor position
type clearScreen struct{}

// Clear the current screen line. Does not move the cursor
type clearLine struct{}

// Close the screen
type closeScreen struct{}

// Representation of the screen
type screenImpl struct {
	currentCursor cursorPos
	savedCursor   cursorPos
	ch            chan screenCmd
	stdOut        io.Writer
}

const (
	esc            string = "\x1B"
	csi            string = esc + "["
	cursorPosition string = csi + "%d;%dH"
	clrLine        string = csi + "2K"
	clrScreen      string = csi + "2J"

	chanBufferSize = 5
)

var (
	screen *screenImpl
)

// Command implementations

func safeWrite(w io.Writer, s string) int {
	n, err := io.WriteString(w, s)
	if err != nil {
		panic("screen write error")
	}
	return n
}

func (wm cursorPos) writeOn(thisScreen *screenImpl) {
	// Update the cursor position
	thisScreen.currentCursor = cursorPos{row: wm.row, col: wm.col}
	safeWrite(thisScreen.stdOut, fmt.Sprintf(cursorPosition, wm.row, wm.col))
}

func (wm writeScreenLine) writeOn(thisScreen *screenImpl) {
	// Move and clear, write the string, then restore the current position
	str := fmt.Sprintf(
		cursorPosition+clrLine+"%s"+cursorPosition,
		wm.cursor.row, wm.cursor.col,
		wm.str,
		thisScreen.currentCursor.row, thisScreen.currentCursor.col)
	safeWrite(thisScreen.stdOut, str)
}

func (wm writeStr) writeOn(thisScreen *screenImpl) {
	// Update the cursor column
	thisScreen.currentCursor.col += len(wm.str)
	safeWrite(thisScreen.stdOut, wm.str)
}

func (wm clearScreen) writeOn(thisScreen *screenImpl) {
	thisScreen.currentCursor = cursorPos{1, 1}
	safeWrite(thisScreen.stdOut, clrScreen)
}

func (wm clearLine) writeOn(thisScreen *screenImpl) {
	safeWrite(thisScreen.stdOut, clrLine)
}

func (wm closeScreen) writeOn(_ *screenImpl) {}

func WriteScreenLine(row int, col int, s string) {
	screen.ch <- writeScreenLine{
		cursor: cursorPos{
			row: row,
			col: col,
		},
		str: s,
	}
}

func Write(s string) {
	screen.ch <- writeStr{str: s}
}

func ClearScreen() {
	screen.ch <- clearScreen{}
}

func ClearLine() {
	screen.ch <- clearLine{}
}

func PositionCursor(row, col int) {
	screen.ch <- cursorPos{row, col}
}

func Close() {
	screen.ch <- closeScreen{}
}

func Initialize() {
	screen = &screenImpl{currentCursor: cursorPos{1, 1}, ch: make(chan screenCmd, chanBufferSize), stdOut: os.Stdout}
	go func() {
		for {
			msg := <-screen.ch
			msg.writeOn(screen)
		}
	}()

	ClearScreen()
}
