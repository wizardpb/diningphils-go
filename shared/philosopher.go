package shared

import "github.com/wizardpb/diningphils-go/shared/philstate"

// Philosopher is an interface implemented by all algorithms to execute the philosopher behavior
type Philosopher interface {
	GetID() int
	GetState() philstate.Enum
	Messages() chan Message
	Execute(m Message)
	Runnable() bool
	Start()
}

// CreateParams are parametrs for philosopher creation/initialization
type CreateParams struct {
	ID         int
	Name       string
	ThinkRange TimeRange
	EatRange   TimeRange
}

// Factory is a factory function type for creating Forks and Philosophers
type Factory func(params CreateParams) (Philosopher, Fork)

// Run sets up the core run loop - repeatedly receive and execute Messages, then Start. The Start initializes the
// initial state and then sets the Philosopher thinking (generally by calling the base Start() method).
func Run(p Philosopher) {
	go func() {
		for m := range p.Messages() {
			p.Execute(m)
		}
	}()
	p.Start()
}
