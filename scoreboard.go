package scoreboard

import "sync"

type Scoreboard struct {
	matches map[string]Match
	nextSeq uint64
	mu      sync.RWMutex
}

func NewScoreboard() *Scoreboard {
	return &Scoreboard{
		matches: make(map[string]Match),
	}
}