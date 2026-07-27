// Package models contains the data structures for API responses from the Highlightly API.
package models

import "encoding/json"

// FlexString is a string that also accepts JSON numbers and null when
// unmarshalling. Statistics values mix types — the same array can contain
// "55%", 12, and null. Numbers keep their literal form (12 becomes "12");
// null becomes the empty string.
type FlexString string

// UnmarshalJSON accepts a JSON string, number, or null.
func (f *FlexString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*f = ""
		return nil
	}
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*f = FlexString(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*f = FlexString(n.String())
	return nil
}

// Pagination represents the paging information in list responses
type Pagination struct {
	TotalCount int `json:"totalCount"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
}

// PlanMeta describes the API plan returned alongside list responses. Message
// carries plan-coverage notices such as "Some results might be hidden with
// FREE tier."
type PlanMeta struct {
	Tier    string `json:"tier"`
	Message string `json:"message"`
}

// Country represents a country
type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

// LeagueSeason is one season a league covers
type LeagueSeason struct {
	Season int `json:"season"`
}

// League represents a league. Top-level league resources carry the Seasons
// list; leagues nested inside matches carry the single Season year instead.
type League struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Logo    string         `json:"logo"`
	Season  int            `json:"season"`
	Seasons []LeagueSeason `json:"seasons"`
	Country Country        `json:"country"`
}

// Team represents a team
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
	Type string `json:"type"`
}

// Player represents a player as returned by the players list
type Player struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`
	Logo     string `json:"logo"`
}
