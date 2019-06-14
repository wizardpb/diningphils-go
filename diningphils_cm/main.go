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
	inCh      chan message
	outCh     chan message
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

// Cycle phil ids 0-4
func nextPhilId(i int) int {
	return (i + 1) % NPhils
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
	}, "Checking dirty on non-existent fork!")
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

	actOn := func(sp *PhilState) {
		switch {
		case p.isHungry() && p.hasForkRequest(sp) && !p.hasFork(sp):
			sp.outCh <- forkRequest{forkId: sp.forkId}
			sp.requested = false
		case !p.isEating() && p.hasForkRequest(sp) && p.isDirty(sp):
			fork := p.forks[sp.forkId]
			delete(p.forks, sp.forkId)
			fork.dirty = false
			sp.outCh <- sentFork{fork: fork}
		default:
			//Do nothing
		}
	}

	if sp != nil {
		actOn(sp)
	} else {
		actOn(p.leftState)
		actOn(p.rightState)
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
		case m = <-p.leftState.inCh:
			m.process(p, p.leftState)
			p.newState(p.leftState)
		case m = <-p.rightState.inCh:
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
	Chans        [NPhils][2]chan message
)

func NewPhilosopher(id int, state int, leftFork *Fork, rightFork *Fork, leftReq bool, rightReq bool) Philosopher {
	var leftIndex, rightIndex = id, nextPhilId(id)

	var leftState, rightState = new(PhilState), new(PhilState)
	leftState.outCh = Chans[leftIndex][0]
	leftState.inCh = Chans[leftIndex][1]
	leftState.forkId = leftIndex
	leftState.requested = leftReq
	rightState.outCh = Chans[rightIndex][1]
	rightState.inCh = Chans[rightIndex][0]
	rightState.forkId = rightIndex
	rightState.requested = rightReq

	forks := make(map[int]*Fork)
	if leftFork != nil {
		forks[leftIndex] = leftFork
	}
	if rightFork != nil {
		forks[rightIndex] = rightFork
	}

	return Philosopher{
		runFlag: true,
		id:      id, name: PhilNames[id],
		state:     state,
		forks:     forks,
		tickChan:  make(chan message, 1),
		leftState: leftState, rightState: rightState,
	}
}

func setup() {

	for i := range PhilNames {
		Forks[i] = Fork{id: i, dirty: true}
		Chans[i][0], Chans[i][1] = make(chan message, 1), make(chan message, 1)
	}

	for i := range PhilNames {
		leftFork, rightFork := &Forks[i], &Forks[nextPhilId(i)]
		leftReq, rightReq := false, false

		/*
		 * Set up forks so the dependency graph is acyclic: phil 0 has both forks, phil 1 has none,
		 * the rest have the left fork only
		 *
		 * Request flags are set opposite this so that all philospohers can initial request the missing fork
		 */
		switch i {
		case 0:
			break // Already set
		case 1:
			leftFork, rightFork = nil, nil
			leftReq, rightReq = true, true
		default:
			rightFork = nil
			rightReq = true
		}

		Philosophers[i] = NewPhilosopher(i, THINKING, leftFork, rightFork, leftReq, rightReq)
	}
}

func main() {

	setup()

	for _, p := range Philosophers {
		p.tickChan <- runState{false}
	}
}
