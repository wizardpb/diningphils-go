package fingers

import (
	"github.com/wizardpb/diningphils-go/shared"
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

type Philosopher struct {
	*shared.PhilosopherBase
}

func (p *Philosopher) NewState() {
	switch p.PhilosopherBase.State {
	case philstate.Hungry:
		// When eating w
		p.PhilosopherBase.State = philstate.Eating
		p.PhilosopherBase.Eat()
	case philstate.Thinking:
		p.PhilosopherBase.Think()
	}
}

func Factory(params shared.CreateParams) (shared.Philosopher, shared.Fork) {
	return &Philosopher{&shared.PhilosopherBase{
			ID:          params.ID,
			Name:        params.Name,
			State:       philstate.Inactive,
			ThinkRange:  params.ThinkRange,
			EatRange:    params.EatRange,
			MessageChan: make(chan shared.Message, 0),
		}}, &shared.ForkBase{
			ID:    params.ID,
			Owner: -1,
		}
}
