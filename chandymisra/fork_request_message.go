package chandymisra

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/shared"
)

// ForkRequestMessage requests a ForkMessage from a Philosopher. Reception causes the fork request flag to be set
type ForkRequestMessage struct {
	Requester shared.Philosopher
	Fork      shared.Fork
}

func (m ForkRequestMessage) String() string {
	return fmt.Sprintf("Philosopher %d requests fork %d", asPhilosopher(m.Requester).ID, asFork(m.Fork).ID)
}

func (m ForkRequestMessage) Process(p shared.Philosopher) bool {
	mp := asPhilosopher(p)
	shared.Assert(func() bool { return !mp.HasRequestFor(m.Fork) }, "fork %d has already been requested", m.Fork.GetID())
	mp.WriteString(fmt.Sprintf("received fork request for %d", m.Fork.GetID()))
	mp.SetRequested(m.Fork, true)
	return true
}
