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
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// PlanMeta describes the API plan and quota usage returned alongside list responses
type PlanMeta struct {
	Tier         string `json:"tier"`
	RequestLimit int    `json:"requestLimit"`
	RequestCount int    `json:"requestCount"`
}

// Country represents a country
type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

// League represents a league for a given season
type League struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Logo    string  `json:"logo"`
	Season  int     `json:"season"`
	Country Country `json:"country"`
}

// Team represents a team
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

// Player represents a player
type Player struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Number   int    `json:"number"`
	Position string `json:"position"`
}

// StatisticValue is a single named statistic
type StatisticValue struct {
	Name  string     `json:"name"`
	Value FlexString `json:"value"`
}
