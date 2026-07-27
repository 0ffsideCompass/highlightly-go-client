package client

import (
	"encoding/json"
	"fmt"

	"github.com/0ffsideCompass/highlightly-go-client/models"
)

const (
	countriesEndpoint = "countries"
	leaguesEndpoint   = "leagues"
)

// Countries hits the /countries endpoint
/*
	- name: (Type: string)
	  Filter countries by name. Value format: England

	All parameters are optional.
*/
func (c *Client) Countries(
	params map[string]any,
) (*models.CountryList, error) {
	endpointURL := fmt.Sprintf("%s%s", c.Domain, countriesEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting countries: %w", err)
	}

	var resp models.CountryList
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling countries response: %w",
			err,
		)
	}

	return &resp, nil
}

// Country hits the /countries/{countryCode} endpoint. It returns a single
// country by its ISO code (for example "GB").
func (c *Client) Country(countryCode string) (*models.Country, error) {
	if countryCode == "" {
		return nil, fmt.Errorf("'countryCode' is required")
	}

	endpointURL := fmt.Sprintf("%s%s/%s", c.Domain, countriesEndpoint, countryCode)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting country: %w", err)
	}

	var resp models.Country
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling country response: %w",
			err,
		)
	}

	return &resp, nil
}

// Leagues hits the /leagues endpoint
/*
	- limit: (Type: integer)
	  The maximum number of leagues to return. Value format: 20
	- offset: (Type: integer)
	  The number of leagues to skip. Value format: 40
	- season: (Type: integer) (should be 4 digits)
	  The year of the season. Value format: 2023
	- leagueName: (Type: string)
	  The name of the league. Value format: Premier League
	- countryCode: (Type: string)
	  The ISO code of the country. Value format: GB
	- countryName: (Type: string)
	  The name of the country. Value format: England

	All parameters are optional.
*/
func (c *Client) Leagues(
	params map[string]any,
) (*models.LeagueList, error) {
	endpointURL := fmt.Sprintf("%s%s", c.Domain, leaguesEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting leagues: %w", err)
	}

	var resp models.LeagueList
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling leagues response: %w",
			err,
		)
	}

	return &resp, nil
}

// League hits the /leagues/{id} endpoint. It returns a single league by ID.
func (c *Client) League(id int) (*models.League, error) {
	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, leaguesEndpoint, id)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting league: %w", err)
	}

	var resp models.League
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling league response: %w",
			err,
		)
	}

	return &resp, nil
}
