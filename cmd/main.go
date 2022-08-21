package main

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/fingers"
	"github.com/wizardpb/diningphils-go/screen"
	"github.com/wizardpb/diningphils-go/shared"
	"time"
)

const (
	NPhils   int = 5
	ThinkMin     = 5
	ThinkMax     = 15
	EatMin       = 5
	EatMax       = 15
)

var (
	PhilNames    [NPhils]string = [5]string{"Arendt", "Butler", "Churchland", "deBeauvoir", "Themistoclea"}
	Philosophers [NPhils]shared.Philosopher
	Forks        [NPhils]*shared.Fork
)

func Initialize(f shared.Factory) {
	// Create Forks and Philosophers
	for i, name := range PhilNames {
		Forks[i] = &shared.Fork{
			ID:    i,
			Owner: -1,
		}

		params := shared.CreateParams{
			ID:         i,
			Name:       name,
			ThinkRange: shared.TimeRange{Min: ThinkMin, Max: ThinkMax, Unit: time.Second},
			EatRange:   shared.TimeRange{Min: EatMin, Max: EatMax, Unit: time.Second},
		}
		Philosophers[i] = f(params)
	}
}

func main() {

	screen.Clear()
	Initialize(fingers.Factory)

	for _, p := range Philosophers {
		shared.Run(p)
	}

	readCmd()
}

const PromptLine = NPhils + 2

func readCmd() {
	screen.Prompt(PromptLine)
	run := true
	for run {
		var cmd string
		_, err := fmt.Scanln(&cmd)
		if err != nil {
			screen.Write(PromptLine, err.Error())
		}
		switch cmd {
		case "q", "Q":
			run = false
		default:
			break
		}
		screen.Prompt(PromptLine)
	}
	screen.Clear()
}
