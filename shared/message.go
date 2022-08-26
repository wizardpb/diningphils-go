package shared

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

type Message interface {
	fmt.Stringer
}

type NewState struct {
	NewState philstate.Enum
}

func (m NewState) String() string {
	return "NewState: " + m.NewState.String()
}
