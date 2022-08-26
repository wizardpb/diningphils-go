package shared

import (
	"fmt"
	"github.com/wizardpb/diningphils-go/screen"
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

// PhilosopherBase implements features commons to all algorithm implementations
type PhilosopherBase struct {
	ID          int
	Name        string
	State       philstate.Enum
	ThinkRange  TimeRange
	EatRange    TimeRange
	MessageChan chan Message
}

// StartThinking - philosopher is thinking, arrange for them to go hungry
func (pb *PhilosopherBase) StartThinking() {
	pb.WriteString("starts thinking")
	pb.DelaySend(pb.ThinkRange, NewState{NewState: philstate.Hungry})
}

// StartEating - set the philosopher eating, arrange for them to finish and think
func (pb *PhilosopherBase) StartEating() {
	pb.WriteString("starts eating")
	pb.DelaySend(pb.EatRange, NewState{NewState: philstate.Thinking})
}

// Eat - set the Philosopher in the Eat state
func (pb *PhilosopherBase) Eat() {

	pb.State = philstate.Eating
	pb.StartEating()
}

// CheckEating checks two invariants that should be true when a Philosopher eats.
// Implement this separately because the 'fingers' implementation intentionally violates this
func (pb *PhilosopherBase) CheckEating() {
	// Check the primary invariant - neither neighbor should be eating, and I should hold
	// both forks
	Assert(
		func() bool {
			return pb.LeftPhilosopher().GetState() != philstate.Eating &&
				pb.RightPhilosopher().GetState() != philstate.Eating
		},
		"eat while a neighbor is eating",
	)

	Assert(
		func() bool {
			return pb.LeftFork().IsHeldBy(pb.ID) &&
				pb.RightFork().IsHeldBy(pb.ID)
		},
		"eat without holding forks",
	)
}

// IsHungry - is the philosopher hungry?
func (pb *PhilosopherBase) IsHungry() bool {
	return pb.State == philstate.Hungry
}

// IsEating - is the philosopher eating?
func (pb *PhilosopherBase) IsEating() bool {
	return pb.State == philstate.Eating
}

// GetID returns the philosopher ID
func (pb *PhilosopherBase) GetID() int {
	return pb.ID
}

func (pb *PhilosopherBase) GetState() philstate.Enum {
	return pb.State
}

// Messages implements part of the Philosopher interface
func (pb *PhilosopherBase) Messages() chan Message {
	return pb.MessageChan
}

// Start implements part of the Philosopher interface
func (pb *PhilosopherBase) Start() {
	pb.State = philstate.Thinking
	pb.StartThinking()
}

// Runnable implements part of the Philosopher interface
func (pb *PhilosopherBase) Runnable() bool {
	return pb.State != philstate.Stopped
}

// WriteString writes a string to the screen on the line dedicated to the philosopher
func (pb *PhilosopherBase) WriteString(s string) {
	forkState := ""
	switch {
	case pb.HoldsFork(pb.LeftFork()) && pb.HoldsFork(pb.RightFork()):
		forkState = fmt.Sprintf(", holds forks %d and %d", pb.leftForkID(), pb.rightForkID())
	case pb.HoldsFork(pb.LeftFork()):
		forkState = fmt.Sprintf(", holds fork %d", pb.leftForkID())
	case pb.HoldsFork(pb.RightFork()):
		forkState = fmt.Sprintf(", holds fork %d", pb.rightForkID())
	}
	screen.Write(pb.ID+1, fmt.Sprintf("%s (%d,%s) %s%s", pb.Name, pb.ID, pb.State, s, forkState))
	//fmt.Println(fmt.Sprintf("%s(%d) %s%s ", pb.Name, pb.ID, s, forkState))
}

// DelaySend sends the given messages to the Philosopher after a random wait given by t
func (pb *PhilosopherBase) DelaySend(t TimeRange, m Message) {
	SendIn(RandDuration(t), m, pb)
}

// left fork ID is always the same as the Philosopher ID
func (pb *PhilosopherBase) leftForkID() int {
	return pb.ID
}

// right fork ID is always the same as the Philosopher ID + 1, wrapping around the table
func (pb *PhilosopherBase) rightForkID() int {
	return (pb.ID + 1) % NPhils
}

// LeftFork returns the fork on the Philosophers left - the one at its ID.
func (pb *PhilosopherBase) LeftFork() Fork {
	return Forks[pb.leftForkID()]
}

// RightFork returns the fork on the Philosophers right - the one at its ID + 1,
// wrapping around the table.
func (pb *PhilosopherBase) RightFork() Fork {
	return Forks[pb.rightForkID()]
}

func (pb *PhilosopherBase) IsLeftFork(f Fork) bool {
	return f == pb.LeftFork()
}

func (pb *PhilosopherBase) IsRightFork(f Fork) bool {
	return f == pb.RightFork()
}

func (pb *PhilosopherBase) HoldsFork(f Fork) bool {
	return f.IsHeldBy(pb.ID)
}

func (pb *PhilosopherBase) LeftPhilosopher() Philosopher {
	// Add NPhils to avoid a negative index
	index := (pb.ID + NPhils - 1) % NPhils
	return Philosophers[index]
}

func (pb *PhilosopherBase) RightPhilosopher() Philosopher {
	index := (pb.ID + 1) % NPhils
	return Philosophers[index]
}
