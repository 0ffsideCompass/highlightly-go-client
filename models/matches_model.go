package models

// MatchScore holds the score of a match
type MatchScore struct {
	Current string `json:"current"`
}

// MatchState describes the live state of a match
type MatchState struct {
	Description string     `json:"description"`
	Clock       int        `json:"clock"`
	Score       MatchScore `json:"score"`
}

// Match represents a match or fixture
type Match struct {
	ID       int        `json:"id"`
	Date     string     `json:"date"`
	League   League     `json:"league"`
	HomeTeam Team       `json:"homeTeam"`
	AwayTeam Team       `json:"awayTeam"`
	State    MatchState `json:"state"`
}

// MatchList is the response from the /matches, /last-five-games and
// /head-2-head endpoints
type MatchList struct {
	Data       []Match    `json:"data"`
	Plan       PlanMeta   `json:"plan"`
	Pagination Pagination `json:"pagination"`
}

// Event represents a live match event
type Event struct {
	Time   int    `json:"time"`
	Type   string `json:"type"`
	Team   Team   `json:"team"`
	Player Player `json:"player"`
}

// EventList is the response from the /events/{id} endpoint
type EventList struct {
	Data []Event `json:"data"`
}

// Lineup represents one team's lineup for a match
type Lineup struct {
	Formation   string   `json:"formation"`
	StartingXI  []Player `json:"startingXI"`
	Substitutes []Player `json:"substitutes"`
}

// Lineups is the response from the /lineups/{matchId} endpoint
type Lineups struct {
	MatchID  int    `json:"matchId"`
	HomeTeam Lineup `json:"homeTeam"`
	AwayTeam Lineup `json:"awayTeam"`
}
