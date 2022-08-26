package fingers

import (
	"github.com/wizardpb/diningphils-go/shared"
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

// Philosopher implementation
type Philosopher struct {
	*shared.PhilosopherBase
}

// Execute implements the Philosopher interface for the Fingers implementation. Just start eating
// when hungry
func (p *Philosopher) Execute(m shared.Message) {
	switch mt := m.(type) {
	case shared.NewState:
		// Update our state value
		p.State = mt.NewState
		switch p.State {
		case philstate.Hungry:
			// When eating with fingers - no need to wait for forks!
			p.Eat()
		case philstate.Thinking:
			p.PhilosopherBase.StartThinking()
		}
	default:
		panic("unknown message: " + m.String())
	}
}

// Factory is the Philosopher and Fork creation function
func Factory(params shared.CreateParams) (shared.Philosopher, shared.Fork) {
	return &Philosopher{&shared.PhilosopherBase{
			ID:          params.ID,
			Name:        params.Name,
			State:       philstate.Inactive,
			ThinkRange:  params.ThinkRange,
			EatRange:    params.EatRange,
			MessageChan: make(chan shared.Message, 0),
		}}, &shared.ForkBase{
			ID:     params.ID,
			Holder: shared.UnOwned,
		}
}
