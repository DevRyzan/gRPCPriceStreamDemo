package main

import (
	"log"
	"sync"
	"time"
)

// LevelUpEvent is emitted when a user levels up.
type LevelUpEvent struct {
	UserID    string
	OldLevel  int
	NewLevel  int
	Timestamp int64
}

// LevelUp runs a goroutine that processes level-up events from a channel.
func LevelUp(levelCh <-chan LevelUpEvent, done <-chan struct{}, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	for {
		select {
		case <-done:
			return
		case e, ok := <-levelCh:
			if !ok {
				return
			}
			time.Sleep(15 * time.Millisecond)
			log.Printf("[levelup] user=%s %d -> %d", e.UserID, e.OldLevel, e.NewLevel)
		}
	}
}

// LevelUpUnbuffered processes from an unbuffered channel.
func LevelUpUnbuffered(levelCh <-chan LevelUpEvent, done <-chan struct{}, wg *sync.WaitGroup) {
	LevelUp(levelCh, done, wg)
}

// LevelUpBuffered processes from a buffered channel (caller chooses cap).
func LevelUpBuffered(levelCh <-chan LevelUpEvent, done <-chan struct{}, wg *sync.WaitGroup) {
	LevelUp(levelCh, done, wg)
}
