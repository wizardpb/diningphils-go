package resourcehierarchy

import (
	"github.com/wizardpb/diningphils-go/shared"
	"sync"
)

type Fork struct {
	shared.ForkBase
	cond *sync.Cond // Conditional to sync access
}
