package models

// MatchScore holds the score of a match
type MatchScore struct {
	Current   string `json:"current"`
	Penalties string `json:"penalties"`
}

// MatchState describes the live state of a match
type MatchState struct {
	Description string     `json:"description"`
	Clock       int        `json:"clock"`
	Score       MatchScore `json:"score"`
}

// Venue describes where a match is played. Only present on the match detail
// resource.
type Venue struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Capacity int    `json:"capacity"`
}

// Referee describes a match official. Only present on the match detail
// resource.
type Referee struct {
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
}

// Forecast describes the weather forecast for a match. Only present on the
// match detail resource.
type Forecast struct {
	Status      string     `json:"status"`
	Temperature FlexString `json:"temperature"`
}

// Match represents a match or fixture. List resources carry the core fields;
// the match detail resource additionally fills Venue, Referee, Forecast and
// Events.
type Match struct {
	ID       int        `json:"id"`
	Round    string     `json:"round"`
	Date     string     `json:"date"`
	Country  Country    `json:"country"`
	League   League     `json:"league"`
	HomeTeam Team       `json:"homeTeam"`
	AwayTeam Team       `json:"awayTeam"`
	State    MatchState `json:"state"`
	Venue    *Venue     `json:"venue,omitempty"`
	Referee  *Referee   `json:"referee,omitempty"`
	Forecast *Forecast  `json:"forecast,omitempty"`
	Events   []Event    `json:"events,omitempty"`
}

// MatchList is the response from the /matches endpoint
type MatchList struct {
	Data       []Match    `json:"data"`
	Plan       PlanMeta   `json:"plan"`
	Pagination Pagination `json:"pagination"`
}

// Event represents a match event (goal, card, substitution, ...). Player and
// Assist are display names; Substituted is the name of the player leaving the
// pitch on a substitution.
type Event struct {
	Team              Team   `json:"team"`
	Time              string `json:"time"`
	Type              string `json:"type"`
	Player            string `json:"player"`
	PlayerID          int    `json:"playerId"`
	Assist            string `json:"assist"`
	AssistingPlayerID int    `json:"assistingPlayerId"`
	Substituted       string `json:"substituted"`
}

// EventList is the response from the /events/{id} endpoint, which returns a
// bare JSON array.
type EventList []Event

// LineupTeam is one team's lineup for a match, with the team fields inlined.
type LineupTeam struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Logo          string   `json:"logo"`
	Formation     string   `json:"formation"`
	InitialLineup []Player `json:"initialLineup"`
	Substitutes   []Player `json:"substitutes"`
}

// Lineups is the response from the /lineups/{matchId} endpoint
type Lineups struct {
	HomeTeam LineupTeam `json:"homeTeam"`
	AwayTeam LineupTeam `json:"awayTeam"`
}
