package client

import (
	"encoding/json"
	"fmt"

	"github.com/0ffsideCompass/highlightly-go-client/models"
)

const (
	highlightsEndpoint                = "highlights"
	highlightsGeoRestrictionsEndpoint = "highlights/geo-restrictions"
)

// Highlights hits the /highlights endpoint
/*
	- countryCode: (Type: string)
	  The ISO code of the country. Value format: GB
	- countryName: (Type: string)
	  The name of the country. Value format: England
	- leagueName: (Type: string)
	  The name of the league. Value format: Premier League
	- leagueId: (Type: integer)
	  The ID of the league. Value format: 2486
	- date: (Type: string) (format: YYYY-MM-DD)
	  The date of the match. Value format: 2023-08-12
	- timezone: (Type: string)
	  The timezone for match dates. Value format: Europe/London
	- season: (Type: integer) (should be 4 digits)
	  The year of the season. Value format: 2023
	- matchId: (Type: integer)
	  The ID of the match. Value format: 1028043
	- homeTeamId: (Type: integer)
	  The ID of the home team. Value format: 3057
	- awayTeamId: (Type: integer)
	  The ID of the away team. Value format: 3068
	- homeTeamName: (Type: string)
	  The name of the home team. Value format: Arsenal
	- awayTeamName: (Type: string)
	  The name of the away team. Value format: Chelsea
	- limit: (Type: integer)
	  The maximum number of highlights to return. Value format: 20
	- offset: (Type: integer)
	  The number of highlights to skip. Value format: 40

	All parameters are optional.
*/
func (c *Client) Highlights(
	params map[string]any,
) (*models.HighlightList, error) {
	endpointURL := fmt.Sprintf("%s%s", c.Domain, highlightsEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting highlights: %w", err)
	}

	var resp models.HighlightList
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling highlights response: %w",
			err,
		)
	}

	return &resp, nil
}

// Highlight hits the /highlights/{id} endpoint. It returns a single highlight
// clip by ID.
func (c *Client) Highlight(id int) (*models.Highlight, error) {
	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, highlightsEndpoint, id)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting highlight: %w", err)
	}

	var resp models.Highlight
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling highlight response: %w",
			err,
		)
	}

	return &resp, nil
}

// HighlightGeoRestrictions hits the /highlights/geo-restrictions/{id}
// endpoint. It returns the countries in which a highlight is blocked or
// allowed.
func (c *Client) HighlightGeoRestrictions(id int) (*models.GeoRestriction, error) {
	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, highlightsGeoRestrictionsEndpoint, id)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting highlight geo restrictions: %w", err)
	}

	var resp models.GeoRestriction
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling highlight geo restrictions response: %w",
			err,
		)
	}

	return &resp, nil
}
