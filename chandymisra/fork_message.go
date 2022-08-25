package chandymisra

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/shared"
)

// ForkMessage is a messages that sends a ForkMessage to a Philosopher
type ForkMessage struct {
	Sender shared.Philosopher
	Fork   shared.Fork
}

func (m ForkMessage) String() string {
	return fmt.Sprintf("Philosopher %d sends fork %d", asPhilosopher(m.Sender).ID, asFork(m.Fork).ID)
}

func (m ForkMessage) Process(p shared.Philosopher) bool {
	mp := asPhilosopher(p)
	mf := asFork(m.Fork)

	mp.WriteString(fmt.Sprintf("receives fork %d", mf.ID))
	shared.Assert(func() bool { return !mf.IsHeld() }, "fork %d already held by %d", mf.ID, mf.Holder)
	mf.SetHolder(mp.ID)

	// If we have both forks we can now eat! Both forks are now dirty
	lf := asFork(mp.LeftFork())
	rf := asFork(mp.RightFork())

	if lf.IsHeldBy(mp.ID) && rf.IsHeldBy(mp.ID) {
		mp.WriteString(fmt.Sprintf("holds both forks and can eat"))
		mp.Eat()
		lf.Dirty = true
		rf.Dirty = true
	}
	return true
}
