package shared

import (
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

type Philosopher interface {
	GetID() int
	SendMessage(m Message)
	Messages() chan Message
	NewState()
	SetState(enum philstate.Enum)
	Runnable() bool
	Start()
}

type CreateParams struct {
	ID         int
	Name       string
	ThinkRange TimeRange
	EatRange   TimeRange
}

type Factory func(params CreateParams) (Philosopher, Fork)

func Run(p Philosopher) {
	go func() {
		for m := range p.Messages() {
			if m.Process(p) {
				p.NewState()
			}
		}
	}()
	p.Start()
}
