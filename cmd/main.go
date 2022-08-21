package main

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/fingers"
	"github.com/wizardpb/diningphils-go/resourcehierarchy"
	"github.com/wizardpb/diningphils-go/screen"
	"github.com/wizardpb/diningphils-go/shared"
	"os"
	"time"
)

func Initialize(f shared.Factory) {
	// Create Forks and Philosophers
	for i, name := range shared.PhilNames {
		params := shared.CreateParams{
			ID:         i,
			Name:       name,
			ThinkRange: shared.TimeRange{Min: shared.ThinkMin, Max: shared.ThinkMax, Unit: time.Second},
			EatRange:   shared.TimeRange{Min: shared.EatMin, Max: shared.EatMax, Unit: time.Second},
		}
		shared.Philosophers[i], shared.Forks[i] = f(params)
	}
}

func writeString(f *os.File, s string) {
	_, err := f.WriteString(s)
	if err != nil {
		os.Exit(4)
	}
}

func main() {

	if len(os.Args) < 2 {
		writeString(os.Stderr, "missing implementation argument")
		os.Exit(1)
	}

	screen.Clear()
	switch os.Args[1] {
	case "fingers", "f":
		Initialize(fingers.Factory)
	case "resourcehierarchy", "rh":
		Initialize(resourcehierarchy.Factory)
	default:
		writeString(os.Stderr, "unknown implementation: "+os.Args[1])
		os.Exit(2)
	}

	for _, p := range shared.Philosophers {
		shared.Run(p)
	}

	readCmd()
}

const PromptLine = shared.NPhils + 2

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
