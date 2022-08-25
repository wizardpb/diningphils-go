package chandymisra

import "github.com/wizardpb/diningphils-go/shared"

type Fork struct {
	*shared.ForkBase
	Dirty bool
}

func asFork(f shared.Fork) *Fork {
	return f.(*Fork)
}
