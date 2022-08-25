package resourcehierarchy

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/shared"
	"github.com/wizardpb/diningphils-go/shared/philstate"
	"sync"
)

type Philosopher struct {
	*shared.PhilosopherBase
	// fork order is a 2-tuple that defines the order the pick-up order of the left and right
	// forks (forkOrder[0] first). It is set by the factory function so that Philosophers 0 to NPhils-2 pick up the
	// right fork first, and Philosopher[NPhils-1[ picks up the left fork first.
	//
	// This enforces the resource hierarchy, and ensure deadlock-free operation
	forkOrder [2]*Fork
}

func (p *Philosopher) NewState() {
	switch p.State {
	case philstate.Hungry:
		for _, f := range p.forkOrder {
			p.pickUp(f)
		}
		p.Eat()
	case philstate.Thinking:
		for _, f := range p.forkOrder {
			p.putDown(f)
		}
		p.Think()
	}
}

func (p *Philosopher) pickUp(f *Fork) {
	f.cond.L.Lock()
	// IsHeld() is correct here because we never pick up a fork without it being put down first
	// The algorithm ensures that pickUp always 'happens before' putDown
	for f.IsHeld() {
		p.WriteString(fmt.Sprintf("is waiting for fork %d", f.ID))
		f.cond.Wait()
	}
	p.WriteString(fmt.Sprintf("has fork %d", f.ID))
	f.SetHolder(p.ID)
	f.cond.L.Unlock()
}

func (p *Philosopher) putDown(f *Fork) {
	f.cond.L.Lock()
	f.SetFree()
	p.WriteString(fmt.Sprintf("puts down fork %d", f.ID))
	f.cond.L.Unlock()
	f.cond.Signal()
}

func Factory(params shared.CreateParams) (shared.Philosopher, shared.Fork) {

	p := &Philosopher{
		PhilosopherBase: &shared.PhilosopherBase{
			ID:          params.ID,
			Name:        params.Name,
			State:       philstate.Inactive,
			ThinkRange:  params.ThinkRange,
			EatRange:    params.EatRange,
			MessageChan: make(chan shared.Message, 0),
		}}

	f := &Fork{
		ForkBase: shared.ForkBase{
			ID:     params.ID,
			Holder: shared.UnOwned,
		},
		cond: sync.NewCond(new(sync.Mutex)),
	}

	return p, f
}

func (p *Philosopher) Start() {
	// Determine fork order
	if p.ID == shared.NPhils-1 {
		// Highest Philosopher picks right first
		p.forkOrder[0] = p.RightFork().(*Fork)
		p.forkOrder[1] = p.LeftFork().(*Fork)
	} else {
		// All others pick up the left first
		p.forkOrder[0] = p.LeftFork().(*Fork)
		p.forkOrder[1] = p.RightFork().(*Fork)
	}

	// Then actually start
	p.PhilosopherBase.Start()
}
