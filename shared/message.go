package shared

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

// Message is an abstract type representing some communication to a Philosopher. Is is use to implement both internal
// state changes (e.g, Thinking to Hungry), and also messages relevant to an algorithm (e.g. fork request messages in
// the C&M implementation)
type Message interface {
	fmt.Stringer
}

// NewState is a message whic cuases a Philosoper to change state (e.g, Thinking to Hungry)
type NewState struct {
	NewState philstate.Enum
}

// String is the implementation of the Stringer interface
func (m NewState) String() string {
	return "NewState: " + m.NewState.String()
}
