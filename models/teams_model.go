package models

// TeamList is the response from the /teams endpoint
type TeamList struct {
	Data       []Team     `json:"data"`
	Plan       PlanMeta   `json:"plan"`
	Pagination Pagination `json:"pagination"`
}
