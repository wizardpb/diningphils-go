package chandymisra

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/shared"
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

// Philosopher implementation
type Philosopher struct {
	*shared.PhilosopherBase
	// ForkRequest holds fork requests for left (index 0) and right (index 1) Forks (nil => not requested).
	ForkRequest [2]bool
}

// Index of the fork request flag of fork f
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

// Convert to my implementation
func asPhilosopher(p shared.Philosopher) *Philosopher {
	return p.(*Philosopher)
}

// Who should I ask for Fork f ?
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

// Set the request flag for Fork f
func (p *Philosopher) setRequested(f shared.Fork, b bool) {
	p.ForkRequest[p.flagID(f)] = b
}

// Do I have a request for Fork f?
func (p *Philosopher) hasRequestFor(f shared.Fork) bool {
	return p.ForkRequest[p.flagID(f)]
}

// Eat - check the invariants and dirty the forks before starting to eat
func (p *Philosopher) Eat() {
	p.CheckEating()
	// Dirty the forks first...
	for _, mf := range []*Fork{asFork(p.LeftFork()), asFork(p.RightFork())} {
		shared.Assert(func() bool { return mf.IsHeldBy(p.ID) }, " eating without holding a fork")
		mf.Dirty = true
	}
	p.PhilosopherBase.Eat()
}

// Execute implements the primary guarded command specified in the C&M paper:
// https://www.cs.utexas.edu/users/misra/scannedPdf.dir/DrinkingPhil.pdf
func (p *Philosopher) Execute(m shared.Message) {
	// Update any state change indicated by the message...
	switch mt := m.(type) {

	case shared.NewState:
		// Update our state value
		p.State = mt.NewState
		switch p.State {
		case philstate.Hungry:
			// No action here - taken care of below
			break
		case philstate.Thinking:
			p.PhilosopherBase.StartThinking()
		}

	case ForkMessage:
		// C&M (R4) - receive a fork
		f := asFork(mt.Fork)
		shared.Assert(func() bool { return !f.IsHeld() }, "fork %d already held by %d", f.ID, f.Holder)
		p.WriteString(fmt.Sprintf("receives fork %d", f.ID))
		f.SetHolder(p.ID)

		// If we have both forks we can now eat! Both forks will now be dirty
		if p.LeftFork().IsHeldBy(p.ID) && p.RightFork().IsHeldBy(p.ID) {
			p.WriteString(fmt.Sprintf("holds both forks and can eat"))
			p.Eat()
		}

	case ForkRequestMessage:
		//C&M (R3) - receive a fork request
		shared.Assert(func() bool { return !p.hasRequestFor(mt.Fork) }, "fork %d has already been requested", mt.Fork.GetID())
		p.WriteString(fmt.Sprintf("received fork request for %d", mt.Fork.GetID()))
		p.setRequested(mt.Fork, true)

	default:
		p.WriteString("unknown message: " + m.String())
	}

	// ... and then check for any implied message send.
	// Fork each of my forks...
	for _, f := range []shared.Fork{p.LeftFork(), p.RightFork()} {
		mf := asFork(f)
		switch {

		case p.IsHungry() && p.hasRequestFor(f) && !f.IsHeldBy(p.ID):
			// C&M (R1): I'm hungry and I need a fork - request it from the appropriate philosopher
			p.setRequested(f, false)
			p.philosopherFor(f).Messages() <- ForkRequestMessage{
				Requester: p,
				Fork:      f,
			}
			p.WriteString(fmt.Sprintf("requested fork %d", f.GetID()))

		case !p.IsEating() && p.hasRequestFor(f) && f.IsHeldBy(p.ID) && mf.Dirty:
			// C&M (R2): I'm done eating and someone has requested a fork - free it (and clean it) then send it over
			requestingPhilosopher := p.philosopherFor(f)
			mf.Dirty = false
			mf.SetFree()
			requestingPhilosopher.Messages() <- ForkMessage{
				Sender: p,
				Fork:   f,
			}
			p.WriteString(fmt.Sprintf("sent fork %d to philosopher %d", f.GetID(), requestingPhilosopher.GetID()))
		}
	}
}

// Factory is the creation function for a Philosopher
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

// Start implements the Philosopher interface
func (p *Philosopher) Start() {
	/*
	 * Set initial conditions:
	 *
	 * Set up forks so the dependency graph is acyclic: phil 0 has both forks, phil 1 has none,
	 * the rest have the left fork only
	 *
	 * Request flags are set opposite this so that all philosophers can initially request the missing fork. The default
	 * value is false so we only need to set the flag for the missing fork
	 */

	for _, p := range shared.Philosophers {
		mp := asPhilosopher(p)
		switch mp.ID {
		case 0:
			// Both request flags false
			mp.LeftFork().SetHolder(mp.ID)
			mp.RightFork().SetHolder(mp.ID)
		case 1:
			// Forks missing, set both request flags
			mp.setRequested(mp.LeftFork(), true)
			mp.setRequested(mp.RightFork(), true)
		default:
			// Hold the left fork, set the right request flag
			mp.LeftFork().SetHolder(mp.ID)
			mp.setRequested(mp.RightFork(), true)
		}
	}

	// Then actually start
	p.PhilosopherBase.Start()
}
