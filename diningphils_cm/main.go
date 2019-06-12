package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

const (
	nanos = 1000000000 // # of nanoseconds in a second

	THINKING = 1
	HUNGRY   = 2
	EATING   = 3

	minEat   = 5
	maxEat   = 20
	minThink = 7
	maxThink = 25

	NPhils = 5
)

type Fork struct {
	id    int
	dirty bool
}

type PhilState struct {
	ch        chan message
	forkId    int
	requested bool
}

type Philosopher struct {
	runFlag               bool
	id                    int
	name                  string
	state                 int
	forks                 map[int]*Fork
	tickChan              chan message
	leftState, rightState *PhilState
}

type message interface {
	process(p *Philosopher, sp *PhilState)
}

// message implementations

// nextState - set the next Philosopher state
type nextState struct {
	state int
}

// forkRequest - indicate a fork request
type forkRequest struct {
	forkId int
}

// sentFork - indicate a requested Fork is being sent
type sentFork struct {
	fork *Fork
}

type runState struct {
	runFlag bool
}

// Utilities
func assert(f func() bool, msg string) {
	if !f() {
		panic(msg)
	}
}

func checkForkId(p *Philosopher, sp *PhilState, fId int, prefix string) {
	assert(
		func() bool { return fId == sp.forkId },
		fmt.Sprintf("%s, philosopher %d (%s): incorrect forkId: %d", prefix, sp.forkId, p.name, fId))
}

// randomSeconds returns a random Duration, rounded to the nearest second,
// between min and max
func randomSeconds(min, max int) time.Duration {
	var maxR = big.NewInt(int64(max - min))
	t, _ := rand.Int(rand.Reader, maxR)
	return time.Duration((int64(min) + t.Int64()) * nanos)
}

// sendIn arranged to send the given message on the given channel after delaying for delay seconds.
// it returns immediately this is done
func sendIn(msg message, c chan message, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		c <- msg
	}()
}

// Philosopher methods

func (p *Philosopher) hasFork(sp *PhilState) bool {
	_, present := p.forks[sp.forkId]
	return present
}

func (p *Philosopher) hasForkRequest(sp *PhilState) bool {
	return sp.requested
}

func (p *Philosopher) isHungry() bool {
	return p.state == HUNGRY
}

func (p *Philosopher) isEating() bool {
	return p.state == EATING
}

func (p *Philosopher) isDirty(sp *PhilState) bool {
	fork, ok := p.forks[sp.forkId]
	assert(func() bool {
		return ok
	}, "Checking dirty on non-existen fork!")
	return fork.dirty
}

func (p *Philosopher) eat() {
	assert(func() bool { return len(p.forks) == 2 }, "Eating without forks!")
	p.state = EATING
	for _, f := range p.forks {
		f.dirty = true
	}
	sendIn(nextState{state: THINKING}, p.tickChan, randomSeconds(minEat, maxEat))
}

func (p *Philosopher) think() {
	p.state = THINKING
	sendIn(nextState{state: HUNGRY}, p.tickChan, randomSeconds(minThink, maxThink))
}

func (p *Philosopher) hungry() {
	p.state = HUNGRY
}

// Message methods
func (m nextState) process(p *Philosopher, sp *PhilState) {
	switch m.state {
	case EATING:
		p.eat()
	case THINKING:
		p.think()
	case HUNGRY:
		p.hungry()
	}
}

func (m forkRequest) process(p *Philosopher, sp *PhilState) {
	checkForkId(p, sp, m.forkId, "Received fork request")
	sp.requested = true
}

func (m sentFork) process(p *Philosopher, sp *PhilState) {
	checkForkId(p, sp, m.fork.id, "Received fork")
	assert(func() bool { return !m.fork.dirty }, "Received dirty fork!")
	assert(
		func() bool { return !p.hasFork(sp) },
		fmt.Sprintf("Received fork: already has fork %d", m.fork.id))
	p.forks[m.fork.id] = m.fork
}

func (m runState) process(p *Philosopher, sp *PhilState) {
	p.runFlag = m.runFlag
}

// newStat implements the full guarded command as specified in Chandy and Misra]

func (p *Philosopher) newState(sp *PhilState) {
	if sp != nil {
		switch {
		case p.isHungry() && p.hasForkRequest(sp) && !p.hasFork(sp):
			sp.ch <- forkRequest{forkId: sp.forkId}
			sp.requested = false
		case !p.isEating() && p.hasForkRequest(sp) && p.isDirty(sp):
			fork := p.forks[sp.forkId]
			delete(p.forks, sp.forkId)
			fork.dirty = false
			sp.ch <- sentFork{fork: fork}
		default:
			//Do nothing
		}
	}

	// Start eating if I am hungry and have both forks
	if p.isHungry() && p.hasFork(p.leftState) && p.hasFork(p.rightState) {
		p.eat()
	}
}

func (p *Philosopher) run() {
	var m message
	for p.runFlag {
		select {
		case m = <-p.leftState.ch:
			m.process(p, p.leftState)
			p.newState(p.leftState)
		case m = <-p.rightState.ch:
			m.process(p, p.rightState)
			p.newState(p.rightState)
		case m = <-p.tickChan:
			m.process(p, nil)
			p.newState(nil)
		}
	}
}

var (
	PhilNames    [NPhils]string = [NPhils]string{"Kant", "Marx", "Hegel", "Spinoza", "Plato"}
	Philosophers [NPhils]Philosopher
	Forks        [NPhils]Fork
)

func main() {

}
