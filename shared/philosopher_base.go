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

// Think - philosopher is thinking, arrange for them to go hungry
func (pb *PhilosopherBase) Think() {
	pb.WriteString("now thinking")
	pb.DelaySend(pb.ThinkRange, &NewState{NewState: philstate.Hungry})
}

// Eat - set the philosopher eating
func (pb *PhilosopherBase) Eat() {
	pb.WriteString("now eating")
	pb.DelaySend(pb.EatRange, &NewState{NewState: philstate.Thinking})
}

// SetState implements part of the Philosopher interface
func (pb *PhilosopherBase) SetState(newState philstate.Enum) {
	pb.State = newState
}

// SendMessage implements part of the Philosopher interface
func (pb *PhilosopherBase) SendMessage(m Message) {
	pb.MessageChan <- m
}

// RecvMessage implements part of the Philosopher interface
func (pb *PhilosopherBase) RecvMessage() Message {
	m := <-pb.MessageChan
	return m
}

func (pb *PhilosopherBase) Messages() chan Message {
	return pb.MessageChan
}

// Start implements part of the Philosopher interface
func (pb *PhilosopherBase) Start() {
	pb.State = philstate.Thinking
	pb.Think()
}

// Runnable implements part of the Philosopher interface
func (pb *PhilosopherBase) Runnable() bool {
	return pb.State != philstate.Stopped
}

// WriteString writes a string to the screen on the line dedicated to the philosopher
func (pb *PhilosopherBase) WriteString(s string) {
	screen.Write(pb.ID+1, fmt.Sprintf("%s(%d) %s", pb.Name, pb.ID, s))
}

// DelaySend sends the given message to the Philosopher after a random wait given by t
func (pb *PhilosopherBase) DelaySend(t TimeRange, m Message) {
	SendIn(RandDuration(t), m, pb)
}

// LeftFork returns the fork on the Philosophers left - the one at its ID.
func (pb *PhilosopherBase) LeftFork() Fork {
	return Forks[pb.ID]
}

// RightFork returns the fork on the Philosophers right - the one at its ID + 1,
// wrapping around the table.
func (pb *PhilosopherBase) RightFork() Fork {
	return Forks[(pb.ID+1)%NPhils]
}
