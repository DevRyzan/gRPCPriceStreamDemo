package main

import (
	"log"
	"sync"
	"time"
)

// MatchResult holds the outcome of a match for one user.
type MatchResult struct {
	MatchID   string
	UserID    string
	Won       bool
	Score     int
	Duration  time.Duration
	Rewards   []Reward
	NewLevel  int
	LeveledUp bool
}

// MatchEnd is the main method: it processes a match end in a goroutine, then
// can fan out to give rewards, save history, and level up (each could be
// separate channels). This version runs one goroutine per match end that does
// the full flow locally (no integration with the other methods).
func MatchEnd(result MatchResult) {
	go func() {
		log.Printf("[match_end] match=%s user=%s won=%v score=%d", result.MatchID, result.UserID, result.Won, result.Score)
		// Simulate sequential steps
		time.Sleep(20 * time.Millisecond)
		// Here you would send to rewardCh, historyCh, levelCh if integrated
		log.Printf("[match_end] done match=%s", result.MatchID)
	}()
}

// MatchEndUnbuffered processes match ends one at a time: caller sends on an
// unbuffered channel and blocks until the worker receives.
func MatchEndUnbuffered(matchCh <-chan MatchResult, done <-chan struct{}, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	for {
		select {
		case <-done:
			return
		case m, ok := <-matchCh:
			if !ok {
				return
			}
			log.Printf("[match_end] match=%s user=%s won=%v", m.MatchID, m.UserID, m.Won)
			time.Sleep(20 * time.Millisecond)
		}
	}
}

// MatchEndBuffered processes match ends from a buffered channel; caller can
// send multiple matches without blocking until the buffer is full.
func MatchEndBuffered(matchCh <-chan MatchResult, done <-chan struct{}, wg *sync.WaitGroup) {
	MatchEndUnbuffered(matchCh, done, wg)
}
