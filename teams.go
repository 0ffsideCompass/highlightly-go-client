package client

import (
	"encoding/json"
	"fmt"

	"github.com/0ffsideCompass/highlightly-go-client/models"
)

const (
	teamsEndpoint = "teams"
)

// Teams hits the /teams endpoint
/*
	- name: (Type: string)
	  Filter teams by name. Value format: Arsenal
	- type: (Type: string)
	  The type of team. Value format: club
	- limit: (Type: integer)
	  The maximum number of teams to return. Value format: 20
	- offset: (Type: integer)
	  The number of teams to skip. Value format: 40

	All parameters are optional.
*/
func (c *Client) Teams(
	params map[string]any,
) (*models.TeamList, error) {
	endpointURL := fmt.Sprintf("%s%s", c.Domain, teamsEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting teams: %w", err)
	}

	var resp models.TeamList
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling teams response: %w",
			err,
		)
	}

	return &resp, nil
}

// Team hits the /teams/{id} endpoint. It returns a single team by ID.
func (c *Client) Team(id int) (*models.Team, error) {
	endpointURL := fmt.Sprintf("%s%s/%d", c.Domain, teamsEndpoint, id)
	body, err := c.get(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("error getting team: %w", err)
	}

	return firstFromArray[models.Team](body, "team")
}
