package models

// StandingRow is one team's row in a standings table
type StandingRow struct {
	Position int  `json:"position"`
	Team     Team `json:"team"`
	Points   int  `json:"points"`
	Played   int  `json:"played"`
	Won      int  `json:"won"`
	Drawn    int  `json:"drawn"`
	Lost     int  `json:"lost"`
}

// StandingGroup is a named group (league table, conference, ...) of standings rows
type StandingGroup struct {
	Name      string        `json:"name"`
	Standings []StandingRow `json:"standings"`
}

// Standings is the response from the /standings endpoint
type Standings struct {
	League League          `json:"league"`
	Groups []StandingGroup `json:"groups"`
}
