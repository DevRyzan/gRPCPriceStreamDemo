package main

import (
	"sync"
	"testing"
	"time"
)

func TestGiveRewards_Unbuffered(t *testing.T) {
	ch := make(chan Reward)
	done := make(chan struct{})
	var wg sync.WaitGroup
	go GiveRewardsUnbuffered(ch, done, &wg)

	ch <- Reward{UserID: "u1", Amount: 100, Type: "gold"}
	ch <- Reward{UserID: "u2", Amount: 50, Type: "silver"}
	close(ch)
	wg.Wait()
	close(done)
}

func TestGiveRewards_Buffered(t *testing.T) {
	ch := make(chan Reward, 10)
	done := make(chan struct{})
	var wg sync.WaitGroup
	go GiveRewardsBuffered(ch, 10, done, &wg)

	ch <- Reward{UserID: "u1", Amount: 100, Type: "gold"}
	ch <- Reward{UserID: "u2", Amount: 50, Type: "silver"}
	close(ch)
	wg.Wait()
	close(done)
}

func TestSaveHistory_Unbuffered(t *testing.T) {
	ch := make(chan HistoryEntry)
	done := make(chan struct{})
	var wg sync.WaitGroup
	go SaveHistoryUnbuffered(ch, done, &wg)

	ch <- HistoryEntry{UserID: "u1", MatchID: "m1", Action: "match_end"}
	close(ch)
	wg.Wait()
	close(done)
}

func TestSaveHistory_Buffered(t *testing.T) {
	ch := make(chan HistoryEntry, 5)
	done := make(chan struct{})
	var wg sync.WaitGroup
	go SaveHistoryBuffered(ch, done, &wg)

	for i := 0; i < 3; i++ {
		ch <- HistoryEntry{UserID: "u1", MatchID: "m1", Action: "match_end"}
	}
	close(ch)
	wg.Wait()
	close(done)
}

func TestLevelUp_Unbuffered(t *testing.T) {
	ch := make(chan LevelUpEvent)
	done := make(chan struct{})
	var wg sync.WaitGroup
	go LevelUpUnbuffered(ch, done, &wg)

	ch <- LevelUpEvent{UserID: "u1", OldLevel: 1, NewLevel: 2, Timestamp: time.Now().Unix()}
	close(ch)
	wg.Wait()
	close(done)
}

func TestLevelUp_Buffered(t *testing.T) {
	ch := make(chan LevelUpEvent, 5)
	done := make(chan struct{})
	var wg sync.WaitGroup
	go LevelUpBuffered(ch, done, &wg)

	ch <- LevelUpEvent{UserID: "u1", OldLevel: 1, NewLevel: 2, Timestamp: time.Now().Unix()}
	close(ch)
	wg.Wait()
	close(done)
}

func TestMatchEnd_FireAndForget(t *testing.T) {
	MatchEnd(MatchResult{MatchID: "m1", UserID: "u1", Won: true, Score: 100})
	time.Sleep(50 * time.Millisecond)
}

func TestMatchEnd_Unbuffered(t *testing.T) {
	ch := make(chan MatchResult)
	done := make(chan struct{})
	var wg sync.WaitGroup
	go MatchEndUnbuffered(ch, done, &wg)

	ch <- MatchResult{MatchID: "m1", UserID: "u1", Won: true, Score: 100}
	close(ch)
	wg.Wait()
	close(done)
}

func TestMatchEnd_Buffered(t *testing.T) {
	ch := make(chan MatchResult, 3)
	done := make(chan struct{})
	var wg sync.WaitGroup
	go MatchEndBuffered(ch, done, &wg)

	ch <- MatchResult{MatchID: "m1", UserID: "u1", Won: true, Score: 100}
	ch <- MatchResult{MatchID: "m2", UserID: "u2", Won: false, Score: 50}
	close(ch)
	wg.Wait()
	close(done)
}
