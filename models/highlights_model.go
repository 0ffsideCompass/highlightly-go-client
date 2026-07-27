package models

// Highlight represents a video highlight clip
type Highlight struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	EmbedURL    string `json:"embedUrl"`
	ImgURL      string `json:"imgUrl"`
	Channel     string `json:"channel"`
	Source      string `json:"source"`
	Category    string `json:"category"`
	Match       Match  `json:"match"`
}

// HighlightList is the response from the /highlights endpoint
type HighlightList struct {
	Data       []Highlight `json:"data"`
	Plan       PlanMeta    `json:"plan"`
	Pagination Pagination  `json:"pagination"`
}

// GeoRestriction is the response from the /highlights/geo-restrictions/{id}
// endpoint. The endpoint is only available on paid plans, so this shape
// follows the published documentation.
type GeoRestriction struct {
	HighlightID      int      `json:"highlightId"`
	BlockedCountries []string `json:"blockedCountries"`
	AllowedCountries []string `json:"allowedCountries"`
}
