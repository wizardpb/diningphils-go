package chandymisra

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/shared"
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

type Philosopher struct {
	*shared.PhilosopherBase
	// ForkRequest holds fork requests for left (index 0) and right (index 1) Forks (nil => not requested).
	ForkRequest [2]bool
}

func (p *Philosopher) flagID(f shared.Fork) int {
	switch {
	case p.LeftFork() == f:
		return 0
	case p.RightFork() == f:
		return 1
	default:
		panic("fork not left or right for flagID")
	}
}

func asPhilosopher(p shared.Philosopher) *Philosopher {
	return p.(*Philosopher)
}

func (p *Philosopher) newStateFor(f shared.Fork) {
	mf := asFork(f)
	switch {
	case p.IsHungry() && p.HasRequestFor(f) && !f.IsHeldBy(p.ID):
		// I'm hungry and I need a fork - request it from the appropriate philosopher
		p.SetRequested(f, false)
		p.philosopherFor(f).Messages() <- ForkRequestMessage{
			Requester: p,
			Fork:      f,
		}
		p.WriteString(fmt.Sprintf("requested fork %d", f.GetID()))
	case !p.IsEating() && p.HasRequestFor(f) && f.IsHeldBy(p.ID) && mf.Dirty:
		// I'm done eating and someone has requested a fork - free it then send
		requestingPhilosopher := p.philosopherFor(f)
		mf.Dirty = false
		mf.SetFree()
		requestingPhilosopher.Messages() <- ForkMessage{
			Sender: p,
			Fork:   f,
		}
		p.WriteString(fmt.Sprintf("sent fork %d to philosopher %d", f.GetID(), requestingPhilosopher.GetID()))
	default:
		// Otherwise, do nothing
	}
}

func (p *Philosopher) philosopherFor(f shared.Fork) shared.Philosopher {
	switch {
	case p.IsLeftFork(f):
		return p.LeftPhilosopher()
	case p.IsRightFork(f):
		return p.RightPhilosopher()
	default:
		panic(fmt.Sprintf("Incorrect fork %d for philosopher %d", f.GetID(), p.ID))
	}
}

func (p *Philosopher) SetRequested(f shared.Fork, b bool) {
	p.ForkRequest[p.flagID(f)] = b
}

func (p *Philosopher) HasRequestFor(f shared.Fork) bool {
	return p.ForkRequest[p.flagID(f)]
}

func (p *Philosopher) SetState(s philstate.Enum) {
	p.WriteString(fmt.Sprintf("sets state %s", s))
	p.PhilosopherBase.SetState(s)
	if s == philstate.Thinking {
		p.Think()
	}
}

func (p *Philosopher) NewState() {
	for _, tp := range []shared.Fork{p.LeftFork(), p.RightFork()} {
		p.newStateFor(tp)
	}

	// Start eating if I am hungry and have both forks
	if p.IsHungry() && p.LeftFork().IsHeldBy(p.ID) && p.RightFork().IsHeldBy(p.ID) {
		p.Eat()
	}
}

func (p *Philosopher) Eat() {
	// Dirty the forks first. We also need to set the state, as the base method doesn't
	for _, mf := range []*Fork{asFork(p.LeftFork()), asFork(p.RightFork())} {
		shared.Assert(func() bool { return mf.IsHeldBy(p.ID) }, " eating without holding a fork")
		mf.Dirty = true
	}
	p.PhilosopherBase.State = philstate.Eating
	p.PhilosopherBase.Eat()
}

func Factory(params shared.CreateParams) (shared.Philosopher, shared.Fork) {
	return &Philosopher{
			PhilosopherBase: &shared.PhilosopherBase{
				ID:         params.ID,
				Name:       params.Name,
				State:      philstate.Inactive,
				ThinkRange: params.ThinkRange,
				EatRange:   params.EatRange,
				// We need a buffered channel here...
				MessageChan: make(chan shared.Message, 10),
			},
			ForkRequest: [2]bool{false, false},
		},
		&Fork{
			ForkBase: &shared.ForkBase{
				ID:     params.ID,
				Holder: shared.UnOwned,
			},
			Dirty: true, // All Forks start out dirty
		}
}

func (p *Philosopher) Start() {
	/*
	 * Set initial conditions:
	 *
	 * Set up forks so the dependency graph is acyclic: phil 0 has both forks, phil 1 has none,
	 * the rest have the left fork only
	 *
	 * Request flags are set opposite this so that all philosophers can initially request the missing fork. Default is false
	 */

	for _, p := range shared.Philosophers {
		mp := asPhilosopher(p)
		switch mp.ID {
		case 0:
			mp.LeftFork().SetHolder(mp.ID)
			mp.RightFork().SetHolder(mp.ID)
		case 1:
			mp.SetRequested(mp.LeftFork(), true)
			mp.SetRequested(mp.RightFork(), true)
		default:
			mp.LeftFork().SetHolder(mp.ID)
			mp.SetRequested(mp.RightFork(), true)
		}
	}

	// Then actually start
	p.PhilosopherBase.Start()
}
