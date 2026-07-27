package models

// Bookmaker represents a bookmaker
type Bookmaker struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// BookmakerList is the response from the /bookmakers endpoint
type BookmakerList struct {
	Data       []Bookmaker `json:"data"`
	Plan       PlanMeta    `json:"plan"`
	Pagination Pagination  `json:"pagination"`
}

// OddsValue is a single labelled odds value within a market
type OddsValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// OddsMarket is a betting market with its odds values
type OddsMarket struct {
	Name string      `json:"name"`
	Odds []OddsValue `json:"odds"`
}

// Odds represents the odds a bookmaker offers for a match
type Odds struct {
	Match     Match        `json:"match"`
	Bookmaker Bookmaker    `json:"bookmaker"`
	Markets   []OddsMarket `json:"markets"`
}

// OddsList is the response from the /odds endpoint
type OddsList struct {
	Data       []Odds     `json:"data"`
	Plan       PlanMeta   `json:"plan"`
	Pagination Pagination `json:"pagination"`
}
