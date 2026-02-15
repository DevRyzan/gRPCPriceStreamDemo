package main

import (
	"log"
	"sync"
	"time"
)

// Reward represents a single reward to give.
type Reward struct {
	UserID string
	Amount int
	Type   string
}

// GiveRewards runs a goroutine that consumes reward requests from a channel and
// "gives" them (e.g. logs or writes elsewhere).
func GiveRewards(rewardCh <-chan Reward, done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		case r, ok := <-rewardCh:
			if !ok {
				return
			}
			// Simulate work
			time.Sleep(10 * time.Millisecond)
			log.Printf("[rewards] user=%s amount=%d type=%s", r.UserID, r.Amount, r.Type)
		}
	}
}

// GiveRewardsUnbuffered starts a goroutine that reads from an unbuffered channel.
// The sender blocks until the worker receives.
func GiveRewardsUnbuffered(rewardCh <-chan Reward, done <-chan struct{}, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	GiveRewards(rewardCh, done)
}

// GiveRewardsBuffered starts a goroutine that reads from a buffered channel (e.g. cap 100).
// The sender blocks only when the buffer is full.
func GiveRewardsBuffered(rewardCh <-chan Reward, bufferSize int, done <-chan struct{}, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	GiveRewards(rewardCh, done)
}
