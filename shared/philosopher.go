package shared

import "github.com/wizardpb/diningphils-go/shared/philstate"

type Philosopher interface {
	GetID() int
	GetState() philstate.Enum
	Messages() chan Message
	Execute(m Message)
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
			p.Execute(m)
		}
	}()
	p.Start()
}
