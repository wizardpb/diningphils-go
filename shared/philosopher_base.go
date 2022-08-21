package shared

import (
	"github.com/wizardpb/diningphils-go/screen"
	"github.com/wizardpb/diningphils-go/shared/philstate"
)

type PhilosopherBase struct {
	ID          int
	Name        string
	State       philstate.Enum
	ThinkRange  TimeRange
	EatRange    TimeRange
	MessageChan chan Message
}

// Thinking - philosopher is thinking, arrange for them to go hungry
func (pb *PhilosopherBase) Thinking() {
	pb.WriteString(pb.Name + " now thinking")
	pb.DelaySend(pb.ThinkRange, &NewState{NewState: philstate.Hungry})
}

// Eat - set the philosopher eating
func (pb *PhilosopherBase) Eat() {
	pb.WriteString(pb.Name + " now eating")
	pb.DelaySend(pb.EatRange, &NewState{NewState: philstate.Thinking})
}

func (pb *PhilosopherBase) SetState(newState philstate.Enum) {
	pb.State = newState
}

func (pb *PhilosopherBase) SendMessage(m Message) {
	pb.MessageChan <- m
}

func (pb *PhilosopherBase) RecvMessage() Message {
	m := <-pb.MessageChan
	return m
}

func (pb *PhilosopherBase) Start() {
	pb.State = philstate.Thinking
	pb.Thinking()
}

func (pb *PhilosopherBase) Runnable() bool {
	return pb.State != philstate.Stopped
}

func (pb *PhilosopherBase) WriteString(s string) {
	screen.Write(pb.ID+1, s)
}

func (pb *PhilosopherBase) DelaySend(t TimeRange, m Message) {
	SendIn(RandDuration(t), m, pb)
}
