package shared

import (
	"crypto/rand"
	"math/big"
	"time"
)

// TimeRange represent a range of times in a given unit
type TimeRange struct {
	Min, Max int
	Unit     time.Duration
}

// RandDuration returns a random duration between the given TimeRange
func RandDuration(r TimeRange) time.Duration {
	delta := big.NewInt(int64(r.Max - r.Min))
	t, _ := rand.Int(rand.Reader, delta)
	return time.Duration(int64(r.Min)+t.Int64()) * r.Unit
}

// SendIn sends the Message m to the Philosopher pb after delay Duration
func SendIn(delay time.Duration, m Message, pb *PhilosopherBase) {
	go func() {
		time.Sleep(delay)
		pb.MessageChan <- m
	}()
}
