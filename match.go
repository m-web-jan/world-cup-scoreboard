package scoreboard

type Match struct {
	HomeTeam  string
	AwayTeam  string
	HomeScore int
	AwayScore int
	sequence  uint64
}

