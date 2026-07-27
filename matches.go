package client

import (
	"encoding/json"
	"fmt"

	"github.com/0ffsideCompass/highlightly-go-client/models"
)

const (
	matchesEndpoint       = "matches"
	lastFiveGamesEndpoint = "last-five-games"
	headToHeadEndpoint    = "head-2-head"
	lineupsEndpoint       = "lineups"
	eventsEndpoint        = "events"
)

// Matches hits the /matches endpoint
/*
	- leagueName: (Type: string)
	  The name of the league. Value format: Premier League
	- leagueId: (Type: integer)
	  The ID of the league. Value format: 2486
	- date: (Type: string) (format: YYYY-MM-DD)
	  The date of the matches. Value format: 2023-08-12
	- timezone: (Type: string)
	  The timezone for match dates. Value format: Europe/London
	- season: (Type: integer) (should be 4 digits)
	  The year of the season. Value format: 2023
	- homeTeamId: (Type: integer)
	  The ID of the home team. Value format: 3057
	- awayTeamId: (Type: integer)
	  The ID of the away team. Value format: 3068
	- homeTeamName: (Type: string)
	  The name of the home team. Value format: Arsenal
	- awayTeamName: (Type: string)
	  The name of the away team. Value format: Chelsea
	- countryCode: (Type: string)
	  The ISO code of the country. Value format: GB
	- countryName: (Type: string)
	  The name of the country. Value format: England
	- limit: (Type: integer)
	  The maximum number of matches to return. Value format: 20
	- offset: (Type: integer)
	  The number of matches to skip. Value format: 40

	All parameters are optional.
*/
func (c *Client) Matches(
	params map[string]any,
) (*models.MatchList, error) {
	endpointURL := fmt.Sprintf("%s%s", c.Domain, matchesEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting matches: %w", err)
	}

	var resp models.MatchList
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling matches response: %w",
			err,
		)
	}

	return &resp, nil
}

// Match hits the /matches/{id} endpoint. It returns a single match by ID.
func (c *Client) Match(id int) (*models.Match, error) {
	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, matchesEndpoint, id)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting match: %w", err)
	}

	return firstFromArray[models.Match](body, "match")
}

// LastFiveGames hits the /last-five-games endpoint. The endpoint returns a
// bare JSON array of matches.
/*
	- teamId: (Type: integer) (Required)
	  The ID of the team. Value format: 3057
*/
func (c *Client) LastFiveGames(
	params map[string]any,
) ([]models.Match, error) {
	// Validate the parameters
	if err := requireIntParams(params, "teamId"); err != nil {
		return nil, err
	}

	endpointURL := fmt.Sprintf("%s%s", c.Domain, lastFiveGamesEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting last five games: %w", err)
	}

	var resp []models.Match
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling last five games response: %w",
			err,
		)
	}

	return resp, nil
}

// HeadToHead hits the /head-2-head endpoint. The endpoint returns a bare
// JSON array of matches.
/*
	- teamIdOne: (Type: integer) (Required)
	  The ID of the first team. Value format: 3057
	- teamIdTwo: (Type: integer) (Required)
	  The ID of the second team. Value format: 3068
*/
func (c *Client) HeadToHead(
	params map[string]any,
) ([]models.Match, error) {
	// Validate the parameters
	if err := requireIntParams(params, "teamIdOne", "teamIdTwo"); err != nil {
		return nil, err
	}

	endpointURL := fmt.Sprintf("%s%s", c.Domain, headToHeadEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting head to head: %w", err)
	}

	var resp []models.Match
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling head to head response: %w",
			err,
		)
	}

	return resp, nil
}

// Lineups hits the /lineups/{matchId} endpoint. It returns both teams'
// lineups for a match.
func (c *Client) Lineups(matchID int) (*models.Lineups, error) {
	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, lineupsEndpoint, matchID)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting lineups: %w", err)
	}

	var resp models.Lineups
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling lineups response: %w",
			err,
		)
	}

	return &resp, nil
}

// Events hits the /events/{id} endpoint. It returns the live events of a
// match.
func (c *Client) Events(matchID int) (*models.EventList, error) {
	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, eventsEndpoint, matchID)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting events: %w", err)
	}

	var resp models.EventList
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling events response: %w",
			err,
		)
	}

	return &resp, nil
}
