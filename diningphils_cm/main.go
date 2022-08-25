package main

import (
	"crypto/rand"
	"fmt"
	"github.com/wizardpb/diningphils-go/screen"
	"log"
	"math/big"
	"strings"
	"time"
)

const (
	nanos = 1000000000 // # of nanoseconds in a second

	INACTIVE = 0
	THINKING = 1
	HUNGRY   = 2
	EATING   = 3

	minEat   = 1
	maxEat   = 5
	minThink = 1
	maxThink = 5

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

// messages implementations

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

// Utilities
func assert(f func() bool, msg string) {
	if !f() {
		log.Panic(msg)
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
	log.Printf("%s (%d) is now eating", p.name, p.id)
	p.state = EATING
	for _, f := range p.forks {
		f.dirty = true
	}
	sendIn(nextState{state: THINKING}, p.tickChan, randomSeconds(minEat, maxEat))
}

func (p *Philosopher) think() {
	log.Printf("%s (%d) is now thinking", p.name, p.id)
	p.state = THINKING
	sendIn(nextState{state: HUNGRY}, p.tickChan, randomSeconds(minThink, maxThink))
}

func (p *Philosopher) hungry() {
	log.Printf("%s (%d) is now hungry", p.name, p.id)
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
	log.Printf("%s (%d) receives request for fork %d", p.name, p.id, sp.forkId)
	sp.requested = true
}

func (m sentFork) process(p *Philosopher, sp *PhilState) {
	checkForkId(p, sp, m.fork.id, "Received fork")
	assert(func() bool { return !m.fork.dirty }, "Received dirty fork!")
	assert(
		func() bool { return !p.hasFork(sp) },
		fmt.Sprintf("Received fork: already has fork %d", m.fork.id))
	log.Printf("%s (%d) receives fork %d", p.name, p.id, sp.forkId)
	p.forks[m.fork.id] = m.fork
}

// newStat implements the full guarded command as specified in Chandy and Misra]

func (p *Philosopher) newState() {

	if p.state == INACTIVE {
		return
	}

	actOn := func(sp *PhilState) {
		switch {
		case p.isHungry() && p.hasForkRequest(sp) && !p.hasFork(sp):
			log.Printf("%s requesting fork %d", p.name, sp.forkId)
			sp.outCh <- forkRequest{forkId: sp.forkId}
			sp.requested = false
		case !p.isEating() && p.hasForkRequest(sp) && p.hasFork(sp) && p.isDirty(sp):
			log.Printf("%s sending fork %d", p.name, sp.forkId)
			fork := p.forks[sp.forkId]
			delete(p.forks, sp.forkId)
			fork.dirty = false
			sp.outCh <- sentFork{fork: fork}
		default:
			//Do nothing
		}
	}

	actOn(p.leftState)
	actOn(p.rightState)

	// Start eating if I am hungry and have both forks
	if p.isHungry() && p.hasFork(p.leftState) && p.hasFork(p.rightState) {
		p.eat()
	}
	log.Printf("%s (%d) now %s, has %s, has request for %s", p.name, p.id, p.eatState(), p.forkState(), p.requestState())
}

func (p *Philosopher) eatState() string {
	stateStr := ""
	switch p.state {
	case INACTIVE:
		stateStr = "Inactive"
	case THINKING:
		stateStr = "Thinking"
	case HUNGRY:
		stateStr = "Hungry"
	case EATING:
		stateStr = "Eating"
	}
	return stateStr
}

func (p *Philosopher) forkState() string {
	s := []string{}
	for _, f := range p.forks {
		s = append(s, fmt.Sprintf("%d", f.id))
	}

	forkStr := strings.Join(s, ",")
	switch len(s) {
	case 0:
		forkStr = "no forks"
	case 1:
		forkStr = "fork " + forkStr
	case 2:
		forkStr = "forks " + forkStr
	}
	return forkStr
}

func (p *Philosopher) requestState() string {
	reqStr := ""
	switch {
	case p.leftState.requested && p.rightState.requested:
		reqStr = fmt.Sprintf("left and right forks (%d and %d)", p.leftState.forkId, p.rightState.forkId)
	case p.leftState.requested && !p.rightState.requested:
		reqStr = fmt.Sprintf("left fork (%d)", p.leftState.forkId)
	case !p.leftState.requested && p.rightState.requested:
		reqStr = fmt.Sprintf("right fork (%d)", p.rightState.forkId)
	case !p.leftState.requested && !p.rightState.requested:
		reqStr = "no forks"
	}
	return reqStr
}

func (p *Philosopher) stateString() string {
	return fmt.Sprintf("%8s: %s, has %s, has request for %s", p.name, p.eatState(), p.forkState(), p.requestState())
}

func (p *Philosopher) run() {
	var m message
	log.Printf("%s now running", p.name)
	for p.runFlag {
		screen.Write(p.id+1, p.stateString())
		select {
		case m = <-p.leftState.inCh:
			m.process(p, p.leftState)
		case m = <-p.rightState.inCh:
			m.process(p, p.rightState)
		case m = <-p.tickChan:
			m.process(p, nil)
		}
		p.newState()

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

		Philosophers[i] = NewPhilosopher(i, INACTIVE, leftFork, rightFork, leftReq, rightReq)
	}
}

func readCmd() {
	var cmd string = ""
	for cmd != "quit" {
		_, err := fmt.Scanln(&cmd)
		if err != nil {
			cmd = "quit"
		}
		screen.Prompt(NPhils + 2)
	}
	screen.Clear()
}

func main() {

	setup()
	screen.Clear()

	for i := range Philosophers {
		p := &Philosophers[i]
		log.Printf("about to run %d, name=%s", p.id, p.name)
		go p.run()
		p.think()
	}

	time.Sleep(time.Millisecond)
	readCmd()
}
