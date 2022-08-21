package shared

import (
	"crypto/rand"
	"math/big"
	"time"
)

type TimeRange struct {
	Min, Max int
	Unit     time.Duration
}

func RandDuration(r TimeRange) time.Duration {
	delta := big.NewInt(int64(r.Max - r.Min))
	t, _ := rand.Int(rand.Reader, delta)
	return time.Duration(int64(r.Min)+t.Int64()) * r.Unit
}

func SendIn(delay time.Duration, m Message, pb *PhilosopherBase) {
	go func() {
		time.Sleep(delay)
		pb.MessageChan <- m
	}()
}
