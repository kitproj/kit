package util

import (
	"runtime"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Semaphores struct {
	values *sync.Map
	seats  map[string]int
}

func NewSemaphores(seatPerKey map[string]int) *Semaphores {
	return &Semaphores{
		values: &sync.Map{},
		seats:  seatPerKey,
	}
}

func (s Semaphores) Get(key string) *semaphore.Weighted {
	seats, ok := s.seats[key]
	if !ok {
		seats = runtime.NumCPU()
	}
	actual, _ := s.values.LoadOrStore(key, semaphore.NewWeighted(int64(seats)))
	mutex := actual.(*semaphore.Weighted)
	return mutex
}
