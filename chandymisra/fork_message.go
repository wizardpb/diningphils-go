package chandymisra

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/shared"
)

// ForkMessage is a messages that sends a Fork to a Philosopher
type ForkMessage struct {
	Sender shared.Philosopher
	Fork   shared.Fork
}

// String implements the Stringer interface
func (m ForkMessage) String() string {
	return fmt.Sprintf("Philosopher %d sends fork %d", asPhilosopher(m.Sender).ID, asFork(m.Fork).ID)
}
