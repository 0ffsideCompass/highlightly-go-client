package client

import (
	"encoding/json"
	"fmt"

	"github.com/0ffsideCompass/highlightly-go-client/models"
)

const (
	teamStatisticsEndpoint  = "teams/statistics"
	matchStatisticsEndpoint = "statistics"
	playersEndpoint         = "players"
)

// TeamStatistics hits the /teams/statistics/{id} endpoint
/*
	- fromDate: (Type: string) (Required) (format: YYYY-MM-DD)
	  The date from which statistics are aggregated. Value format: 2023-08-12
	- timezone: (Type: string)
	  The timezone for match dates. Value format: Europe/London
*/
func (c *Client) TeamStatistics(
	teamID int,
	params map[string]any,
) (*models.TeamStatistics, error) {
	// Validate the parameters
	if err := requireStringParams(params, "fromDate"); err != nil {
		return nil, err
	}

	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, teamStatisticsEndpoint, teamID)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting team statistics: %w", err)
	}

	var resp models.TeamStatistics
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling team statistics response: %w",
			err,
		)
	}

	return &resp, nil
}

// MatchStatistics hits the /statistics/{matchId} endpoint. It returns both
// teams' statistics for a match.
func (c *Client) MatchStatistics(matchID int) (*models.MatchStatistics, error) {
	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, matchStatisticsEndpoint, matchID)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting match statistics: %w", err)
	}

	var resp models.MatchStatistics
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling match statistics response: %w",
			err,
		)
	}

	return &resp, nil
}

// Players hits the /players endpoint
/*
	- name: (Type: string)
	  Filter players by name. Value format: Saka
	- limit: (Type: integer)
	  The maximum number of players to return. Value format: 20
	- offset: (Type: integer)
	  The number of players to skip. Value format: 40

	All parameters are optional.
*/
func (c *Client) Players(
	params map[string]any,
) (*models.PlayerList, error) {
	endpointURL := fmt.Sprintf("%s%s", c.Domain, playersEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting players: %w", err)
	}

	var resp models.PlayerList
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling players response: %w",
			err,
		)
	}

	return &resp, nil
}

// Player hits the /players/{id} endpoint. It returns a single player summary
// by ID.
func (c *Client) Player(id int) (*models.Player, error) {
	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, playersEndpoint, id)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting player: %w", err)
	}

	var resp models.Player
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling player response: %w",
			err,
		)
	}

	return &resp, nil
}

// PlayerStatistics hits the /players/{id}/statistics endpoint. It returns the
// statistics of a player.
func (c *Client) PlayerStatistics(id int) (*models.PlayerStatistics, error) {
	endpointURL := fmt.Sprintf("%s%s/%d/statistics", c.Domain, playersEndpoint, id)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting player statistics: %w", err)
	}

	var resp models.PlayerStatistics
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling player statistics response: %w",
			err,
		)
	}

	return &resp, nil
}
