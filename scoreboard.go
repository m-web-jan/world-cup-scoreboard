package scoreboard

import (
	"slices"
	"sync"
)

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

func (s *Scoreboard) GetSummary() []Match {
	var matches []Match

	s.mu.RLock()
	for _, match := range s.matches {
		matches = append(matches, match)
	}
	s.mu.RUnlock()

	slices.SortFunc(matches, func(a, b Match) int {
		totalA := a.HomeScore + a.AwayScore
		totalB := b.HomeScore + b.AwayScore

		if totalA > totalB {
			return -1
		}
		if totalA < totalB {
			return 1
		}

		if a.sequence > b.sequence {
			return -1
		}
		if a.sequence < b.sequence {
			return 1
		}

		return 0
	})

	return matches
}
