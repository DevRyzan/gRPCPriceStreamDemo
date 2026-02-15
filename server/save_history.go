package main

import (
	"log"
	"sync"
	"time"
)

// HistoryEntry is one row to save (e.g. match result, reward granted).
type HistoryEntry struct {
	UserID    string
	MatchID   string
	Action    string
	Payload   string
	CreatedAt int64
}

// SaveHistory runs a goroutine that reads history entries from a channel and
// "saves" them (e.g. to DB or slice). Use unbuffered for backpressure;
// use buffered to absorb bursts.
func SaveHistory(historyCh <-chan HistoryEntry, done <-chan struct{}, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	for {
		select {
		case <-done:
			return
		case e, ok := <-historyCh:
			if !ok {
				return
			}
			time.Sleep(5 * time.Millisecond)
			log.Printf("[history] user=%s match=%s action=%s", e.UserID, e.MatchID, e.Action)
		}
	}
}

// SaveHistoryUnbuffered expects an unbuffered channel: the producer blocks until
// this goroutine receives. Good when you want strict ordering and backpressure.
func SaveHistoryUnbuffered(historyCh <-chan HistoryEntry, done <-chan struct{}, wg *sync.WaitGroup) {
	SaveHistory(historyCh, done, wg)
}

// SaveHistoryBuffered expects a buffered channel (caller creates with cap > 0).
// Producer can send up to cap items without blocking.
func SaveHistoryBuffered(historyCh <-chan HistoryEntry, done <-chan struct{}, wg *sync.WaitGroup) {
	SaveHistory(historyCh, done, wg)
}
