package client

import (
	"encoding/json"
	"fmt"

	"github.com/0ffsideCompass/highlightly-go-client/models"
)

const (
	bookmakersEndpoint = "bookmakers"
	oddsEndpoint       = "odds"
)

// Bookmakers hits the /bookmakers endpoint
/*
	- name: (Type: string)
	  Filter bookmakers by name. Value format: bet365
	- limit: (Type: integer)
	  The maximum number of bookmakers to return. Value format: 20
	- offset: (Type: integer)
	  The number of bookmakers to skip. Value format: 40

	All parameters are optional.
*/
func (c *Client) Bookmakers(
	params map[string]any,
) (*models.BookmakerList, error) {
	endpointURL := fmt.Sprintf("%s%s", c.Domain, bookmakersEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting bookmakers: %w", err)
	}

	var resp models.BookmakerList
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling bookmakers response: %w",
			err,
		)
	}

	return &resp, nil
}

// Bookmaker hits the /bookmakers/{id} endpoint. It returns a single bookmaker
// by ID.
func (c *Client) Bookmaker(id int) (*models.Bookmaker, error) {
	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, bookmakersEndpoint, id)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting bookmaker: %w", err)
	}

	var resp models.Bookmaker
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling bookmaker response: %w",
			err,
		)
	}

	return &resp, nil
}

// Odds hits the /odds endpoint
/*
	- oddsType: (Type: string)
	  The type of odds. Value format: prematch
	- leagueId: (Type: integer)
	  The ID of the league. Value format: 2486
	- leagueName: (Type: string)
	  The name of the league. Value format: Premier League
	- timezone: (Type: string)
	  The timezone for match dates. Value format: Europe/London
	- bookmakerId: (Type: integer)
	  The ID of the bookmaker. Value format: 1
	- bookmakerName: (Type: string)
	  The name of the bookmaker. Value format: bet365
	- matchId: (Type: integer)
	  The ID of the match. Value format: 1028043
	- date: (Type: string) (format: YYYY-MM-DD)
	  The date of the matches. Value format: 2023-08-12
	- limit: (Type: integer)
	  The maximum number of odds entries to return. Value format: 20
	- offset: (Type: integer)
	  The number of odds entries to skip. Value format: 40

	All parameters are optional.
*/
func (c *Client) Odds(
	params map[string]any,
) (*models.OddsList, error) {
	endpointURL := fmt.Sprintf("%s%s", c.Domain, oddsEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting odds: %w", err)
	}

	var resp models.OddsList
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling odds response: %w",
			err,
		)
	}

	return &resp, nil
}
