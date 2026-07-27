package client_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	client "github.com/0ffsideCompass/highlightly-go-client"
)

func TestMatches(t *testing.T) {
	tests := []struct {
		name          string
		params        map[string]any
		responseBody  string
		statusCode    int
		expectError   bool
		errorContains string
	}{
		{
			name:   "successful matches request",
			params: map[string]any{"leagueId": 2486, "season": 2023},
			responseBody: `{
				"data": [
					{
						"id": 1028043,
						"date": "2023-08-12T14:00:00Z",
						"league": {"id": 2486, "name": "Premier League", "season": 2023},
						"homeTeam": {"id": 3057, "name": "Arsenal"},
						"awayTeam": {"id": 3068, "name": "Chelsea"},
						"state": {"description": "Finished", "clock": 90, "score": {"current": "2-1"}}
					}
				],
				"plan": {"tier": "PRO", "requestLimit": 5000, "requestCount": 12},
				"pagination": {"limit": 20, "offset": 0, "total": 1}
			}`,
			statusCode: http.StatusOK,
		},
		{
			name:          "matches request with API error",
			params:        map[string]any{"leagueId": -1},
			responseBody:  `{"error": "Invalid league"}`,
			statusCode:    http.StatusBadRequest,
			expectError:   true,
			errorContains: "API request failed with status 400",
		},
		{
			name:          "matches request with invalid JSON response",
			params:        map[string]any{"leagueId": 2486},
			responseBody:  `{"invalid": json}`,
			statusCode:    http.StatusOK,
			expectError:   true,
			errorContains: "error unmarshalling matches response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockClient(tt.statusCode, tt.responseBody)
			cli, err := client.New("test-api-key", mock)
			assert.NoError(t, err)

			result, err := cli.Matches(tt.params)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Data, 1)
				assert.Equal(t, 1028043, result.Data[0].ID)
				assert.Equal(t, "Arsenal", result.Data[0].HomeTeam.Name)
				assert.Equal(t, "2-1", result.Data[0].State.Score.Current)
			}
		})
	}
}

func TestLastFiveGamesValidation(t *testing.T) {
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `{}`))
	assert.NoError(t, err)

	result, err := cli.LastFiveGames(map[string]any{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'teamId' is required")
	assert.Nil(t, result)

	result, err = cli.LastFiveGames(map[string]any{"teamId": "not-a-number"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'teamId' must be an integer")
	assert.Nil(t, result)
}

func TestHeadToHeadValidation(t *testing.T) {
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `{}`))
	assert.NoError(t, err)

	result, err := cli.HeadToHead(map[string]any{"teamIdOne": 3057})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'teamIdTwo' is required")
	assert.Nil(t, result)
}

func TestStandings(t *testing.T) {
	t.Run("missing required params", func(t *testing.T) {
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `{}`))
		assert.NoError(t, err)

		result, err := cli.Standings(map[string]any{"season": 2023})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'leagueId' is required")
		assert.Nil(t, result)
	})

	t.Run("successful standings request", func(t *testing.T) {
		responseBody := `{
			"league": {"id": 2486, "name": "Premier League", "season": 2023},
			"groups": [
				{
					"name": "Table",
					"standings": [
						{"position": 1, "team": {"id": 3057, "name": "Arsenal"}, "points": 84, "played": 38, "won": 26, "drawn": 6, "lost": 6}
					]
				}
			]
		}`
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
		assert.NoError(t, err)

		result, err := cli.Standings(map[string]any{"leagueId": 2486, "season": 2023})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Groups, 1)
		assert.Equal(t, 1, result.Groups[0].Standings[0].Position)
		assert.Equal(t, "Arsenal", result.Groups[0].Standings[0].Team.Name)
	})
}

func TestTeamStatistics(t *testing.T) {
	t.Run("missing fromDate", func(t *testing.T) {
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `{}`))
		assert.NoError(t, err)

		result, err := cli.TeamStatistics(3057, map[string]any{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'fromDate' is required")
		assert.Nil(t, result)
	})

	t.Run("successful team statistics request", func(t *testing.T) {
		responseBody := `{
			"team": {"id": 3057, "name": "Arsenal"},
			"matchesPlayed": 38,
			"wins": 26,
			"draws": 6,
			"losses": 6,
			"goalsFor": 88,
			"goalsAgainst": 43
		}`
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
		assert.NoError(t, err)

		result, err := cli.TeamStatistics(3057, map[string]any{"fromDate": "2023-08-12"})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 38, result.MatchesPlayed)
		assert.Equal(t, 88, result.GoalsFor)
	})
}

func TestMatchStatistics(t *testing.T) {
	// Statistics values arrive as mixed strings, numbers, and nulls.
	responseBody := `{
		"matchId": 1028043,
		"teams": [
			{
				"team": {"id": 3057, "name": "Arsenal"},
				"statistics": [
					{"name": "Possession", "value": "55%"},
					{"name": "Shots", "value": 12},
					{"name": "Offsides", "value": null}
				]
			}
		]
	}`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.MatchStatistics(1028043)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	stats := result.Teams[0].Statistics
	assert.Equal(t, "55%", string(stats[0].Value))
	assert.Equal(t, "12", string(stats[1].Value))
	assert.Equal(t, "", string(stats[2].Value))
}

func TestHighlights(t *testing.T) {
	responseBody := `{
		"data": [
			{
				"id": 7701,
				"title": "Arsenal vs Chelsea - Full Highlights",
				"type": "VERIFIED",
				"url": "https://cdn.highlightly.net/clip/7701.mp4",
				"imgUrl": "https://cdn.highlightly.net/clip/7701.jpg",
				"match": {"id": 1028043, "homeTeam": {"id": 3057, "name": "Arsenal"}, "awayTeam": {"id": 3068, "name": "Chelsea"}}
			}
		],
		"plan": {"tier": "PRO", "requestLimit": 5000, "requestCount": 13},
		"pagination": {"limit": 20, "offset": 0, "total": 1}
	}`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.Highlights(map[string]any{"matchId": 1028043})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, 7701, result.Data[0].ID)
	assert.Equal(t, 1028043, result.Data[0].Match.ID)
}

func TestHighlightGeoRestrictions(t *testing.T) {
	responseBody := `{
		"highlightId": 7701,
		"blockedCountries": ["GB", "DE"],
		"allowedCountries": []
	}`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.HighlightGeoRestrictions(7701)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, []string{"GB", "DE"}, result.BlockedCountries)
}

func TestCountry(t *testing.T) {
	t.Run("missing country code", func(t *testing.T) {
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `{}`))
		assert.NoError(t, err)

		result, err := cli.Country("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'countryCode' is required")
		assert.Nil(t, result)
	})

	t.Run("successful country request", func(t *testing.T) {
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `{"code": "GB", "name": "United Kingdom"}`))
		assert.NoError(t, err)

		result, err := cli.Country("GB")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "GB", result.Code)
	})
}

func TestLineups(t *testing.T) {
	responseBody := `{
		"matchId": 1028043,
		"homeTeam": {
			"formation": "4-3-3",
			"startingXI": [{"id": 9001, "name": "Bukayo Saka", "number": 7, "position": "RW"}],
			"substitutes": [{"id": 9002, "name": "Leandro Trossard", "number": 19, "position": "LW"}]
		},
		"awayTeam": {"formation": "4-2-3-1", "startingXI": [], "substitutes": []}
	}`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.Lineups(1028043)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "4-3-3", result.HomeTeam.Formation)
	assert.Equal(t, "Bukayo Saka", result.HomeTeam.StartingXI[0].Name)
}

func TestOdds(t *testing.T) {
	responseBody := `{
		"data": [
			{
				"match": {"id": 1028043},
				"bookmaker": {"id": 1, "name": "bet365"},
				"markets": [
					{
						"name": "Match Winner",
						"odds": [
							{"label": "Home", "value": "1.85"},
							{"label": "Draw", "value": "3.60"},
							{"label": "Away", "value": "4.20"}
						]
					}
				]
			}
		],
		"plan": {"tier": "PRO", "requestLimit": 5000, "requestCount": 14},
		"pagination": {"limit": 20, "offset": 0, "total": 1}
	}`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.Odds(map[string]any{"matchId": 1028043})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "bet365", result.Data[0].Bookmaker.Name)
	assert.Equal(t, "1.85", result.Data[0].Markets[0].Odds[0].Value)
}
