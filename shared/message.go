package shared

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

type Message interface {
	fmt.Stringer
	Process(p Philosopher) bool
}

type NewState struct {
	NewState philstate.Enum
}

func (m NewState) Process(p Philosopher) bool {
	p.SetState(m.NewState)
	return true
}

func (m NewState) String() string {
	return "NewState: " + m.NewState.String()
}
