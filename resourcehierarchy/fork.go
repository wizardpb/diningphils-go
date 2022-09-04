package resourcehierarchy

import (
	"github.com/wizardpb/diningphils-go/shared"
)

type freeToken struct{}

// Fork represents a fork available for eacting. It has shard state, and indicates it's ability to be used using a
// semaphore channel - philosophers wanting the fork wait for the token to appear in the channel. Since it is initialized
// to a single item buffer, any philosopher receiving from the channel will block until the current owner released the
// fork by writing the token. Philosophers compete for a fork by reading the channel simultaneous - the channel semantics
// ensure that only one will (atomically) receive the token and grab the fork.
//
// Note that this does not ensure fairness - a Philosopher that thinks very quickly and repeatedly goes hungry can repeatedly
// grab a fork at the expense of a slower one.
type Fork struct {
	shared.ForkBase
	semChan chan freeToken
}
