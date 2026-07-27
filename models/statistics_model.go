package models

// TeamStatistics is the response from the /teams/statistics/{id} endpoint
type TeamStatistics struct {
	Team          Team `json:"team"`
	MatchesPlayed int  `json:"matchesPlayed"`
	Wins          int  `json:"wins"`
	Draws         int  `json:"draws"`
	Losses        int  `json:"losses"`
	GoalsFor      int  `json:"goalsFor"`
	GoalsAgainst  int  `json:"goalsAgainst"`
}

// TeamMatchStatistics holds one team's statistics within a match
type TeamMatchStatistics struct {
	Team       Team             `json:"team"`
	Statistics []StatisticValue `json:"statistics"`
}

// MatchStatistics is the response from the /statistics/{matchId} endpoint
type MatchStatistics struct {
	MatchID int                   `json:"matchId"`
	Teams   []TeamMatchStatistics `json:"teams"`
}

// PlayerList is the response from the /players endpoint
type PlayerList struct {
	Data       []Player   `json:"data"`
	Plan       PlanMeta   `json:"plan"`
	Pagination Pagination `json:"pagination"`
}

// PlayerStatistics is the response from the /players/{id}/statistics endpoint
type PlayerStatistics struct {
	Player     Player           `json:"player"`
	Statistics []StatisticValue `json:"statistics"`
}
