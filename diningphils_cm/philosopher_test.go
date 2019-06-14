package main

import (
	"fmt"
	"testing"
	"time"
)

func TestChanAssignments(t *testing.T) {

	setup()

	// for all philosophers, their right in and out chans should be the next out and in chans

	for i := range Philosophers {
		p1, p2 := &Philosophers[i], &Philosophers[nextPhilId(i)]
		if !(p1.rightState.inCh == p2.leftState.outCh && p1.rightState.outCh == p2.leftState.inCh) {
			t.Errorf("%d channels not correct", p1.id)
		}
	}
}

func TestRequest(t *testing.T) {
	setup()

	p := &Philosophers[1]

	// Set it going
	go func() { p.run() }()

	// Make p hungry
	p.tickChan <- nextState{state: HUNGRY}

	// p should request both forks
	timeOut := false
	var leftMsg, rightMsg message
	for !timeOut {
		select {
		case leftMsg = <-p.leftState.outCh:
			t.Logf("left received: %v\n", leftMsg)
		case rightMsg = <-p.rightState.outCh:
			t.Logf("right received: %v\n", rightMsg)
		case <-time.After(time.Second):
			timeOut = true
		}
	}

	if leftMsg == nil || rightMsg == nil {
		t.Error(fmt.Sprintf("left=%v, right=%v", leftMsg, rightMsg))
	}

	p.runFlag = false
}

func TestSendFork(t *testing.T) {
	setup()

	p := &Philosophers[0]

	// Set it going
	go func() { p.run() }()

	// p should send a fork when requested

	p.rightState.inCh <- forkRequest{forkId: 1}

	var timeOut bool
	var leftMsg, rightMsg message
	select {
	case leftMsg = <-p.leftState.outCh:
		t.Logf("left received: %v\n", leftMsg)
	case rightMsg = <-p.rightState.outCh:
		t.Logf("right received: %v\n", rightMsg)
	case <-time.After(time.Second):
		timeOut = true
	}

	if timeOut {
		t.Errorf("timed out")
		t.Fail()
	}

	forkMsg := rightMsg.(sentFork)
	if !(leftMsg == nil && forkMsg.fork.id == 1 && !forkMsg.fork.dirty) {
		t.Errorf("bad msg: %v", rightMsg)
	}
}

func TestStartsEating(t *testing.T) {

	setup()

	p := &Philosophers[0]

	// Set it going
	go func() { p.run() }()

	// Make p hungry
	p.tickChan <- nextState{state: HUNGRY}

	time.Sleep(time.Millisecond)

	if !p.isEating() {
		t.Error("p not eating after 1ms")
	}
}
