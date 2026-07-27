package client

import (
	"encoding/json"
	"fmt"

	"github.com/0ffsideCompass/highlightly-go-client/models"
)

const (
	standingsEndpoint = "standings"
)

// Standings hits the /standings endpoint
/*
	- leagueId: (Type: integer) (Required)
	  The ID of the league. Value format: 2486
	- season: (Type: integer) (Required) (should be 4 digits)
	  The year of the season. Value format: 2023
*/
func (c *Client) Standings(
	params map[string]any,
) (*models.Standings, error) {
	// Validate the parameters
	if err := requireIntParams(params, "leagueId", "season"); err != nil {
		return nil, err
	}

	endpointURL := fmt.Sprintf("%s%s", c.Domain, standingsEndpoint)
	body, err := c.get(
		c.buildURL(
			endpointURL,
			params,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting standings: %w", err)
	}

	var resp models.Standings
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf(
			"error unmarshalling standings response: %w",
			err,
		)
	}

	return &resp, nil
}
