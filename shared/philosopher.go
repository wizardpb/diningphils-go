package shared

import (
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

type Philosopher interface {
	SendMessage(m Message)
	RecvMessage() Message
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
		for p.Runnable() {
			m := p.RecvMessage()
			if m.Process(p) {
				p.NewState()
			}
		}
	}()
	p.Start()
}
