package chandymisra

import "github.com/wizardpb/diningphils-go/shared"

// Fork is the C&M Fork implementation. Just add the dirty flag
type Fork struct {
	*shared.ForkBase
	Dirty bool
}

// Convert a Fork interface to my implementation
func asFork(f shared.Fork) *Fork {
	return f.(*Fork)
}
