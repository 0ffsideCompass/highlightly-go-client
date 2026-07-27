package models

// StandingRecord is a win/draw/loss record with goal counts, given separately
// for home, away, and overall.
type StandingRecord struct {
	Games         int `json:"games"`
	Wins          int `json:"wins"`
	Draws         int `json:"draws"`
	Loses         int `json:"loses"`
	ScoredGoals   int `json:"scoredGoals"`
	ReceivedGoals int `json:"receivedGoals"`
}

// StandingRow is one team's row in a standings table
type StandingRow struct {
	Position int            `json:"position"`
	Team     Team           `json:"team"`
	Points   int            `json:"points"`
	Home     StandingRecord `json:"home"`
	Away     StandingRecord `json:"away"`
	Total    StandingRecord `json:"total"`
}

// StandingGroup is a named group (league table, conference, ...) of standings rows
type StandingGroup struct {
	Name      string        `json:"name"`
	Standings []StandingRow `json:"standings"`
}

// Standings is the response from the /standings endpoint
type Standings struct {
	Groups []StandingGroup `json:"groups"`
	League League          `json:"league"`
}
