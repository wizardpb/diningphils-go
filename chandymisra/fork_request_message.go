package chandymisra

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/shared"
)

// ForkRequestMessage requests a Fork from a Philosopher. Reception causes the fork request flag to be set
type ForkRequestMessage struct {
	Requester shared.Philosopher
	Fork      shared.Fork
}

// String implements the Stringer interface
func (m ForkRequestMessage) String() string {
	return fmt.Sprintf("Philosopher %d requests fork %d", asPhilosopher(m.Requester).ID, asFork(m.Fork).ID)
}
