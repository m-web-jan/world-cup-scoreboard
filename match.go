package scoreboard

import (
	"errors"
	"strings"
)

type Match struct {
	HomeTeam  string
	AwayTeam  string
	HomeScore int
	AwayScore int
	sequence  uint64
}

func (s *Scoreboard) StartGame(home, away string) error {
	homeNormalized, awayNormalized := normalizeTeams(home, away)

	if homeNormalized == "" || awayNormalized == "" {
		return errors.New("team name must not be empty")
	}

	if homeNormalized == awayNormalized {
		return errors.New("teams should be different")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := matchKey{
		home: homeNormalized,
		away: awayNormalized,
	}

	if team, inMatch := s.teamInMatch(home, away); inMatch {
		return errors.New(team + " team already in match")
	}

	s.nextSeq++
	newMatch := Match{
		HomeTeam:  strings.TrimSpace(home),
		AwayTeam:  strings.TrimSpace(away),
		HomeScore: 0,
		AwayScore: 0,
		sequence:  s.nextSeq,
	}

	s.matches[key] = newMatch

	return nil
}

func (s *Scoreboard) UpdateScore(home, away string, homeScore, awayScore int) error {
	if homeScore < 0 || awayScore < 0 {
		return errors.New("score should not be negative")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	home, away = normalizeTeams(home, away)
	key := matchKey{
		home: home,
		away: away,
	}

	match, ok := s.matches[key]
	if !ok {
		return errors.New("match does not exist")
	}

	match.HomeScore = homeScore
	match.AwayScore = awayScore

	s.matches[key] = match

	return nil
}

func (s *Scoreboard) FinishGame(home, away string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	home, away = normalizeTeams(home, away)
	key := matchKey{
		home: home,
		away: away,
	}

	if _, ok := s.matches[key]; !ok {
		return errors.New("match does not exist")
	}

	delete(s.matches, key)

	return nil
}

func normalizeTeams(home, away string) (string, string) {
	return strings.ToLower(strings.TrimSpace(home)), strings.ToLower(strings.TrimSpace(away))
}

func (s *Scoreboard) teamInMatch(home, away string) (string, bool) {
	home, away = normalizeTeams(home, away)

	for key := range s.matches {
		if key.home == home || key.away == home {
			return home, true
		}
		if key.home == away || key.away == away {
			return away, true
		}
	}

	return "", false
}
