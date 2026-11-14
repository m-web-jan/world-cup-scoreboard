package scoreboard

import "sync"

type matchKey struct {
	home string
	away string
}

type Scoreboard struct {
	matches map[matchKey]Match
	nextSeq uint64
	mu      sync.RWMutex
}

func NewScoreboard() *Scoreboard {
	return &Scoreboard{
		matches: make(map[matchKey]Match),
	}
}