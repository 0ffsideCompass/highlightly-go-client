package models

import "encoding/json"

// StatisticValue is a single named statistic within a match
type StatisticValue struct {
	DisplayName string     `json:"displayName"`
	Value       FlexString `json:"value"`
}

// TeamMatchStatistics holds one team's statistics within a match
type TeamMatchStatistics struct {
	Team       Team             `json:"team"`
	Statistics []StatisticValue `json:"statistics"`
}

// MatchStatistics is the response from the /statistics/{matchId} endpoint,
// which returns a bare JSON array with one entry per team.
type MatchStatistics []TeamMatchStatistics

// GamesRecord counts played games by outcome
type GamesRecord struct {
	Played int `json:"played"`
	Wins   int `json:"wins"`
	Draws  int `json:"draws"`
	Loses  int `json:"loses"`
}

// GoalsRecord counts goals scored and received
type GoalsRecord struct {
	Scored   int `json:"scored"`
	Received int `json:"received"`
}

// TeamRecord combines game and goal counts for one venue split
type TeamRecord struct {
	Games GamesRecord `json:"games"`
	Goals GoalsRecord `json:"goals"`
}

// TeamLeagueStatistics is one league's record within a team's statistics
type TeamLeagueStatistics struct {
	LeagueID   int        `json:"leagueId"`
	LeagueName string     `json:"leagueName"`
	Season     int        `json:"season"`
	Home       TeamRecord `json:"home"`
	Away       TeamRecord `json:"away"`
	Total      TeamRecord `json:"total"`
}

// TeamStatistics is the response from the /teams/statistics/{id} endpoint,
// which returns a bare JSON array with one entry per league the team plays in.
type TeamStatistics []TeamLeagueStatistics

// PlayerList is the response from the /players endpoint
type PlayerList struct {
	Data       []Player   `json:"data"`
	Plan       PlanMeta   `json:"plan"`
	Pagination Pagination `json:"pagination"`
}

// PlayerClub describes a player's current club engagement
type PlayerClub struct {
	Current                 string `json:"current"`
	JoinedAt                string `json:"joinedAt"`
	ContractExpiry          string `json:"contractExpiry"`
	LatestContractExtension string `json:"latestContractExtension"`
}

// PlayerPosition describes a player's main and secondary positions
type PlayerPosition struct {
	Main      string `json:"main"`
	Secondary string `json:"secondary"`
}

// PlayerProfile is the biographical section of a player detail
type PlayerProfile struct {
	FullName    string         `json:"fullName"`
	Club        PlayerClub     `json:"club"`
	Position    PlayerPosition `json:"position"`
	Foot        string         `json:"foot"`
	Height      string         `json:"height"`
	BirthDate   string         `json:"birthDate"`
	DeathDate   string         `json:"deathDate"`
	BirthPlace  string         `json:"birthPlace"`
	Citizenship string         `json:"citizenship"`
}

// PlayerDetail is the response from the /players/{id} endpoint. The rumour,
// news, transfer, market value and injury sections are kept as raw JSON until
// they are needed in typed form.
type PlayerDetail struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	FullName    string          `json:"fullName"`
	Logo        string          `json:"logo"`
	Profile     PlayerProfile   `json:"profile"`
	Rumours     json.RawMessage `json:"rumours"`
	RelatedNews json.RawMessage `json:"relatedNews"`
	Transfers   json.RawMessage `json:"transfers"`
	MarketValue json.RawMessage `json:"marketValue"`
	Injuries    json.RawMessage `json:"injuries"`
}

// CompetitionStatistics is one competition-season line of a player's statistics
type CompetitionStatistics struct {
	Club              string `json:"club"`
	Type              string `json:"type"`
	League            string `json:"league"`
	Season            string `json:"season"`
	GamesPlayed       int    `json:"gamesPlayed"`
	MinutesPlayed     int    `json:"minutesPlayed"`
	Goals             int    `json:"goals"`
	Assists           int    `json:"assists"`
	OwnGoals          int    `json:"ownGoals"`
	PenaltiesScored   int    `json:"penaltiesScored"`
	YellowCards       int    `json:"yellowCards"`
	SecondYellowCards int    `json:"secondYellowCards"`
	RedCards          int    `json:"redCards"`
	CleanSheets       int    `json:"cleanSheets"`
	GoalsConceded     int    `json:"goalsConceded"`
	SubstitutedIn     int    `json:"substitutedIn"`
	SubstitutedOut    int    `json:"substitutedOut"`
}

// PlayerStatistics is the response from the /players/{id}/statistics endpoint
type PlayerStatistics struct {
	ID             int                     `json:"id"`
	Name           string                  `json:"name"`
	FullName       string                  `json:"fullName"`
	Logo           string                  `json:"logo"`
	PerCompetition []CompetitionStatistics `json:"perCompetition"`
}
