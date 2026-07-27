package models

// CountryList is the response from the /countries endpoint
type CountryList struct {
	Data []Country `json:"data"`
}

// LeagueList is the response from the /leagues endpoint
type LeagueList struct {
	Data       []League   `json:"data"`
	Plan       PlanMeta   `json:"plan"`
	Pagination Pagination `json:"pagination"`
}
