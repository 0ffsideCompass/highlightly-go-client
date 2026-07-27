package models

// CountryList is the response from the /countries endpoint, which returns a
// bare JSON array.
type CountryList []Country

// LeagueList is the response from the /leagues endpoint
type LeagueList struct {
	Data       []League   `json:"data"`
	Plan       PlanMeta   `json:"plan"`
	Pagination Pagination `json:"pagination"`
}
