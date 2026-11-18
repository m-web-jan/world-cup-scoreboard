package scoreboard

import (
	"reflect"
	"testing"
)

func TestStartGame(t *testing.T) {
	s := NewScoreboard()

	err := s.StartGame("Mexico", "Canada")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	games := s.GetSummary()
	if len(games) != 1 {
		t.Fatalf("expected 1 game, got %d", len(games))
	}

	game := games[0]
	if game.HomeTeam != "Mexico" || game.AwayTeam != "Canada" {
		t.Errorf("unexpected teams: %s - %s", game.HomeTeam, game.AwayTeam)
	}

	if game.HomeScore != 0 || game.AwayScore != 0 {
		t.Errorf("expected score 0:0, got %d:%d", game.HomeScore, game.AwayScore)
	}
}

func TestStartGameValidation(t *testing.T) {
	tests := []struct {
		name string
		home string
		away string
	}{
		{"empty home", "", "Canada"},
		{"empty away", "Mexico", ""},
		{"whitespace home", "   ", "Canada"},
		{"same teams", "Mexico", "Mexico"},
		{"same teams different case", "Mexico", "mexico"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScoreboard()

			err := s.StartGame(tt.home, tt.away)
			if err == nil {
				t.Errorf("expected error for %q vs %q", tt.home, tt.away)
			}
		})
	}
}

func TestUpdateScore(t *testing.T) {
	s := NewScoreboard()

	if err := s.StartGame("Mexico", "Canada"); err != nil {
		t.Fatal(err)
	}

	if err := s.UpdateScore("Mexico", "Canada", 2, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	game := s.GetSummary()[0]
	if game.HomeScore != 2 || game.AwayScore != 3 {
		t.Errorf("expected 2:3, got %d:%d", game.HomeScore, game.AwayScore)
	}

	if err := s.UpdateScore("Mexico", "Canada", 1, 1); err != nil {
		t.Fatalf("unexpected error while correcting score: %v", err)
	}

	game = s.GetSummary()[0]
	if game.HomeScore != 1 || game.AwayScore != 1 {
		t.Errorf("expected 1:1, got %d:%d", game.HomeScore, game.AwayScore)
	}
}

func TestFinishGame(t *testing.T) {
	s := NewScoreboard()

	if err := s.StartGame("Mexico", "Canada"); err != nil {
		t.Fatal(err)
	}

	if err := s.FinishGame("Mexico", "Canada"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(s.GetSummary()) != 0 {
		t.Error("expected empty scoreboard after finishing game")
	}

	if err := s.StartGame("Mexico", "Brazil"); err != nil {
		t.Errorf("expected team to be available: %v", err)
	}
}

func TestGetSummaryEmpty(t *testing.T) {
	s := NewScoreboard()

	games := s.GetSummary()
	if len(games) != 0 {
		t.Errorf("expected empty summary, got %d games", len(games))
	}
}

func TestGetSummaryOrder(t *testing.T) {
	s := NewScoreboard()

	tests := []struct {
		home      string
		away      string
		homeScore int
		awayScore int
	}{
		{"Mexico", "Canada", 0, 5},
		{"Spain", "Brazil", 10, 2},
		{"Germany", "France", 2, 2},
		{"Uruguay", "Italy", 6, 6},
		{"Argentina", "Australia", 3, 1},
	}

	for _, tt := range tests {
		if err := s.StartGame(tt.home, tt.away); err != nil {
			t.Fatal(err)
		}

		if err := s.UpdateScore(
			tt.home, tt.away, tt.homeScore, tt.awayScore,
		); err != nil {
			t.Fatal(err)
		}
	}

	games := s.GetSummary()

	expected := []string{
		"Uruguay",
		"Spain",
		"Mexico",
		"Argentina",
		"Germany",
	}

	var actual []string
	for _, game := range games {
		actual = append(actual, game.HomeTeam)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected order %v, got %v", expected, actual)
	}
}

func TestGetSummarySameTotal(t *testing.T) {
	s := NewScoreboard()

	if err := s.StartGame("Mexico", "Canada"); err != nil {
		t.Fatal(err)
	}
	if err := s.StartGame("Spain", "Brazil"); err != nil {
		t.Fatal(err)
	}

	if err := s.UpdateScore("Mexico", "Canada", 2, 2); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateScore("Spain", "Brazil", 3, 1); err != nil {
		t.Fatal(err)
	}

	games := s.GetSummary()

	if len(games) != 2 {
		t.Fatalf("expected 2 games, got %d", len(games))
	}

	if games[0].HomeTeam != "Spain" {
		t.Errorf("expected Spain first, got %s", games[0].HomeTeam)
	}

	if err := s.UpdateScore("Mexico", "Canada", 1, 3); err != nil {
		t.Fatal(err)
	}

	games = s.GetSummary()
	if games[0].HomeTeam != "Spain" {
		t.Error("score update should not change insertion order")
	}
}
