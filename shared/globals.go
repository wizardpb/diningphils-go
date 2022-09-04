package shared

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/screen"
)

// Control constants for timings, number of Philosophers,etc.
const (
	NPhils   int = 5
	ThinkMin     = 5
	ThinkMax     = 15
	EatMin       = 5
	EatMax       = 15

	ScreenPos  = 3
	PromptLine = ScreenPos + NPhils + 3

	promptString = "> "
)

// Philosophers are numbers 0 to NPhils-1, as are forks.
//
// The fork to the left of Philosopher[i] is Fork[i]; the fork to the right is Fork[i+1 mod NPhils], since
// the wrap around the table
var (
	PhilNames    [NPhils]string = [NPhils]string{"Hannah Arendt", "Judith Butler", "Patricia Churchland", "Simone de Beauvoir", "Themistoclea"}
	Philosophers [NPhils]Philosopher
	Forks        [NPhils]Fork
)

// ReadCmd reads and executes a command string from the terminal
func ReadCmd() string {
	screen.PositionCursor(PromptLine, 1)
	screen.ClearLine()
	screen.Write(promptString)
	var cmd string
	_, err := fmt.Scanln(&cmd)
	if err != nil {
		panic("screen read error")
	}
	return cmd
}
