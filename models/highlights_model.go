package models

// Highlight represents a video highlight clip
type Highlight struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	URL    string `json:"url"`
	ImgURL string `json:"imgUrl"`
	Match  Match  `json:"match"`
}

// HighlightList is the response from the /highlights endpoint
type HighlightList struct {
	Data       []Highlight `json:"data"`
	Plan       PlanMeta    `json:"plan"`
	Pagination Pagination  `json:"pagination"`
}

// GeoRestriction is the response from the /highlights/geo-restrictions/{id} endpoint
type GeoRestriction struct {
	HighlightID      int      `json:"highlightId"`
	BlockedCountries []string `json:"blockedCountries"`
	AllowedCountries []string `json:"allowedCountries"`
}
