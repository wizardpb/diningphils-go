package main

import (
	"crypto/rand"
	"fmt"
	"github.com/wizardpb/diningphils-go/screen"
	"log"
	"math/big"
	"os"
	"sync"
	"time"
)

type Fork struct {
	id    int
	owner *Philosopher // The current owner, nil if free
	cond  *sync.Cond   // Conditional to sync access
}

type Philosopher struct {
	id                int
	name              string
	lowFork, highFork *Fork
}

const (
	NPhils   int = 5
	ThinkMin     = 5
	ThinkMax     = 15
	EatMin       = 5
	EatMax       = 15
	Nanos        = 1000000000
)

var (
	PhilNames    [NPhils]string = [5]string{"Kant", "Marx", "Hegel", "Spinoza", "Plato"}
	Philosophers [NPhils]Philosopher
	Forks        [NPhils]Fork
)

func setup() {
	for i := range Forks {
		Forks[i] = Fork{id: i + 1, cond: sync.NewCond(new(sync.Mutex))}
	}
	for i, n := range PhilNames {
		leftId, rightId := i, (i+1)%NPhils
		var lowFork, highFork *Fork
		if leftId < rightId {
			lowFork, highFork = &Forks[leftId], &Forks[rightId]
		} else {
			lowFork, highFork = &Forks[rightId], &Forks[leftId]
		}

		Philosophers[i] = Philosopher{id: i + 1, name: n, lowFork: lowFork, highFork: highFork}
	}
	fp, err := os.Create("app.log")
	if err != nil {
		panic(err.Error())
	}
	log.SetOutput(fp)
}

func (p *Philosopher) pickUp(f *Fork) {
	f.cond.L.Lock()
	for f.owner != nil {
		screen.Instance.Write(p.id, fmt.Sprintf("%s waiting for fork %d", p.name, f.id))
		f.cond.Wait()
	}
	f.owner = p
	//log.Printf("Phil %d grabs fork %d", p.id, f.id)
	f.cond.L.Unlock()
}

func (p *Philosopher) putDown(f *Fork) {
	f.cond.L.Lock()
	f.owner = nil
	//log.Printf("Philosopher %d outs down fork %d",p.id,f.id)
	f.cond.L.Unlock()
	f.cond.Signal()
}

func randSecs(min int, max int) int64 {
	var maxR = big.NewInt(int64(max - min))
	t, _ := rand.Int(rand.Reader, maxR)
	return (int64(min) + t.Int64())
}

func (p Philosopher) think() {
	secs := randSecs(ThinkMin, ThinkMax)
	screen.Instance.Write(p.id, fmt.Sprintf("%s thinking for %d secs", p.name, secs))
	//log.Printf("Philosopher %d thinking for %d secs",p.id,secs)
	time.Sleep(time.Duration(secs * Nanos))
}

func (p Philosopher) eat() {
	secs := randSecs(EatMin, EatMax)
	screen.Instance.Write(p.id, fmt.Sprintf("%s eating for %d secs", p.name, secs))
	//log.Printf("Philosopher %d eat for %d secs",p.id,secs)
	time.Sleep(time.Duration(secs * Nanos))
}

func (p Philosopher) run() {
	for {
		p.think()
		p.pickUp(p.lowFork)
		p.pickUp(p.highFork)
		p.eat()
		p.putDown(p.lowFork)
		p.putDown(p.highFork)
	}
}

func readCmd() {
	var cmd string = ""
	for cmd != "quit" {
		_, err := fmt.Scanln(&cmd)
		if err != nil {
			cmd = "quit"
		}
		screen.Instance.Prompt(NPhils + 2)
	}
	screen.Instance.Clear()
}

func main() {
	screen.Instance.Clear()
	screen.Instance.Prompt(NPhils + 2)
	setup()
	for _, p := range Philosophers {
		go p.run()
	}
	readCmd()
}
