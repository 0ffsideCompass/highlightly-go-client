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
			params: map[string]any{"leagueId": 33973, "season": 2025},
			responseBody: `{
				"data": [
					{
						"id": 1299388429,
						"round": "Regular Season - 15",
						"date": "2026-07-26T23:30:00.000Z",
						"country": {"code": "GB-ENG", "name": "England", "logo": "https://highlightly.net/soccer/images/countries/GB-ENG.svg"},
						"league": {"id": 33973, "name": "Premier League", "season": 2025},
						"homeTeam": {"id": 36526, "name": "Arsenal"},
						"awayTeam": {"id": 43334, "name": "Manchester City"},
						"state": {"description": "Finished", "clock": 90, "score": {"current": "2 - 1", "penalties": null}}
					}
				],
				"plan": {"tier": "BASIC", "message": "All data available with current plan."},
				"pagination": {"totalCount": 402, "offset": 0, "limit": 1}
			}`,
			statusCode: http.StatusOK,
		},
		{
			name:          "matches request with API error",
			params:        map[string]any{"leagueId": -1},
			responseBody:  `{"message": "Invalid league", "statusCode": 400}`,
			statusCode:    http.StatusBadRequest,
			expectError:   true,
			errorContains: "API request failed with status 400",
		},
		{
			name:          "matches request with invalid JSON response",
			params:        map[string]any{"leagueId": 33973},
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
				assert.Equal(t, 1299388429, result.Data[0].ID)
				assert.Equal(t, "Regular Season - 15", result.Data[0].Round)
				assert.Equal(t, "Arsenal", result.Data[0].HomeTeam.Name)
				assert.Equal(t, "2 - 1", result.Data[0].State.Score.Current)
				assert.Equal(t, 402, result.Pagination.TotalCount)
				assert.Equal(t, "BASIC", result.Plan.Tier)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	t.Run("single match decodes from array", func(t *testing.T) {
		responseBody := `[
			{
				"id": 1281156605,
				"round": "Regular Season - 16",
				"date": "2026-07-26T19:00:00.000Z",
				"venue": {"city": null, "name": null, "country": null, "capacity": null},
				"referee": {"name": null, "nationality": null},
				"forecast": {"status": null, "temperature": null},
				"events": [
					{"team": {"id": 1989571, "name": "Ñublense"}, "time": "23", "type": "Goal", "assist": "M. Plaza",
					 "player": "S. Valencia", "playerId": 91557473, "substituted": null, "assistingPlayerId": 49908061}
				],
				"homeTeam": {"id": 1989571, "name": "Ñublense"},
				"awayTeam": {"id": 1973402, "name": "Palestino"},
				"league": {"id": 226299, "name": "Primera División", "season": 2026},
				"state": {"clock": 90, "score": {"current": "2 - 0", "penalties": null}, "description": "Finished"}
			}
		]`
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
		assert.NoError(t, err)

		result, err := cli.Match(1281156605)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1281156605, result.ID)
		assert.Len(t, result.Events, 1)
		assert.Equal(t, "Goal", result.Events[0].Type)
		assert.Equal(t, "S. Valencia", result.Events[0].Player)
	})

	t.Run("empty array is not found", func(t *testing.T) {
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `[]`))
		assert.NoError(t, err)

		result, err := cli.Match(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "match not found")
		assert.Nil(t, result)
	})
}

func TestLastFiveGames(t *testing.T) {
	t.Run("validation", func(t *testing.T) {
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `[]`))
		assert.NoError(t, err)

		result, err := cli.LastFiveGames(map[string]any{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'teamId' is required")
		assert.Nil(t, result)

		result, err = cli.LastFiveGames(map[string]any{"teamId": "not-a-number"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'teamId' must be an integer")
		assert.Nil(t, result)
	})

	t.Run("decodes bare array", func(t *testing.T) {
		responseBody := `[
			{"id": 1314260505, "round": "Final",
			 "state": {"clock": 120, "score": {"current": "1 - 1", "penalties": "4 - 3"}, "description": "Finished after penalties"}}
		]`
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
		assert.NoError(t, err)

		result, err := cli.LastFiveGames(map[string]any{"teamId": 36526})
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "4 - 3", result[0].State.Score.Penalties)
	})
}

func TestHeadToHeadValidation(t *testing.T) {
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `[]`))
	assert.NoError(t, err)

	result, err := cli.HeadToHead(map[string]any{"teamIdOne": 36526})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'teamIdTwo' is required")
	assert.Nil(t, result)
}

func TestStandings(t *testing.T) {
	t.Run("missing required params", func(t *testing.T) {
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `{}`))
		assert.NoError(t, err)

		result, err := cli.Standings(map[string]any{"season": 2025})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'leagueId' is required")
		assert.Nil(t, result)
	})

	t.Run("successful standings request", func(t *testing.T) {
		responseBody := `{
			"groups": [
				{
					"name": "Premier League",
					"standings": [
						{
							"position": 1,
							"points": 85,
							"team": {"id": 36526, "name": "Arsenal"},
							"home": {"wins": 15, "draws": 2, "games": 19, "loses": 2, "scoredGoals": 41, "receivedGoals": 11},
							"away": {"wins": 11, "draws": 5, "games": 19, "loses": 3, "scoredGoals": 30, "receivedGoals": 16},
							"total": {"wins": 26, "draws": 7, "games": 38, "loses": 5, "scoredGoals": 71, "receivedGoals": 27}
						}
					]
				}
			],
			"league": {"id": 33973, "name": "Premier League", "season": 2025}
		}`
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
		assert.NoError(t, err)

		result, err := cli.Standings(map[string]any{"leagueId": 33973, "season": 2025})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Groups, 1)
		row := result.Groups[0].Standings[0]
		assert.Equal(t, 1, row.Position)
		assert.Equal(t, 85, row.Points)
		assert.Equal(t, "Arsenal", row.Team.Name)
		assert.Equal(t, 26, row.Total.Wins)
		assert.Equal(t, 41, row.Home.ScoredGoals)
		assert.Equal(t, "Premier League", result.League.Name)
	})
}

func TestTeamStatistics(t *testing.T) {
	t.Run("missing fromDate", func(t *testing.T) {
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `[]`))
		assert.NoError(t, err)

		result, err := cli.TeamStatistics(36526, map[string]any{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'fromDate' is required")
		assert.Nil(t, result)
	})

	t.Run("successful team statistics request", func(t *testing.T) {
		responseBody := `[
			{
				"leagueId": 33973,
				"season": 2026,
				"leagueName": "Premier League",
				"home": {"games": {"played": 19, "wins": 15, "loses": 2, "draws": 2}, "goals": {"scored": 41, "received": 11}},
				"away": {"games": {"played": 19, "wins": 11, "loses": 3, "draws": 5}, "goals": {"scored": 30, "received": 16}},
				"total": {"games": {"played": 38, "wins": 26, "loses": 5, "draws": 7}, "goals": {"scored": 71, "received": 27}}
			}
		]`
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
		assert.NoError(t, err)

		result, err := cli.TeamStatistics(36526, map[string]any{"fromDate": "2025-08-01"})
		assert.NoError(t, err)
		assert.NotNil(t, result)
		stats := *result
		assert.Len(t, stats, 1)
		assert.Equal(t, "Premier League", stats[0].LeagueName)
		assert.Equal(t, 38, stats[0].Total.Games.Played)
		assert.Equal(t, 71, stats[0].Total.Goals.Scored)
	})
}

func TestMatchStatistics(t *testing.T) {
	// Statistics values arrive as numbers; FlexString keeps their literal form.
	responseBody := `[
		{
			"team": {"id": 1989571, "name": "Ñublense"},
			"statistics": [
				{"value": 112, "displayName": "Attacks"},
				{"value": "55%", "displayName": "Possession"},
				{"value": null, "displayName": "Offsides"}
			]
		}
	]`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.MatchStatistics(1281156605)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	teams := *result
	assert.Len(t, teams, 1)
	stats := teams[0].Statistics
	assert.Equal(t, "Attacks", stats[0].DisplayName)
	assert.Equal(t, "112", string(stats[0].Value))
	assert.Equal(t, "55%", string(stats[1].Value))
	assert.Equal(t, "", string(stats[2].Value))
}

func TestHighlights(t *testing.T) {
	responseBody := `{
		"data": [
			{
				"id": 572811,
				"type": "VERIFIED",
				"imgUrl": "https://i.ytimg.com/vi/6g6GL8aldhY/hqdefault.jpg",
				"title": "TNT Sports Replay | Ñublense 2-0 Palestino | Matchday 16",
				"description": null,
				"url": "https://www.youtube.com/watch?v=6g6GL8aldhY",
				"embedUrl": "https://www.youtube.com/embed/6g6GL8aldhY",
				"channel": "TNT Sports Chile",
				"source": "youtube",
				"category": "match-highlights",
				"match": {"id": 1281156605, "homeTeam": {"id": 1989571, "name": "Ñublense"}, "awayTeam": {"id": 1973402, "name": "Palestino"}}
			}
		],
		"plan": {"tier": "BASIC", "message": "All data available with current plan."},
		"pagination": {"totalCount": 31, "offset": 0, "limit": 2}
	}`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.Highlights(map[string]any{"date": "2026-07-26"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, 572811, result.Data[0].ID)
	assert.Equal(t, "https://www.youtube.com/embed/6g6GL8aldhY", result.Data[0].EmbedURL)
	assert.Equal(t, "youtube", result.Data[0].Source)
	assert.Equal(t, "match-highlights", result.Data[0].Category)
	assert.Equal(t, 1281156605, result.Data[0].Match.ID)
	assert.Equal(t, 31, result.Pagination.TotalCount)
}

func TestHighlight(t *testing.T) {
	responseBody := `[
		{
			"id": 572811,
			"type": "VERIFIED",
			"title": "TNT Sports Replay | Ñublense 2-0 Palestino | Matchday 16",
			"url": "https://www.youtube.com/watch?v=6g6GL8aldhY",
			"channel": "TNT Sports Chile",
			"match": {"id": 1281156605}
		}
	]`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.Highlight(572811)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 572811, result.ID)
	assert.Equal(t, "TNT Sports Chile", result.Channel)
}

func TestHighlightGeoRestrictions(t *testing.T) {
	responseBody := `{
		"highlightId": 572811,
		"blockedCountries": ["GB", "DE"],
		"allowedCountries": []
	}`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.HighlightGeoRestrictions(572811)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, []string{"GB", "DE"}, result.BlockedCountries)
}

func TestCountries(t *testing.T) {
	// The countries list is a bare JSON array.
	responseBody := `[{"code": "GB-ENG", "name": "England", "logo": "https://highlightly.net/soccer/images/countries/GB-ENG.svg"}]`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.Countries(map[string]any{"name": "England"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 1)
	assert.Equal(t, "GB-ENG", (*result)[0].Code)
}

func TestCountry(t *testing.T) {
	t.Run("missing country code", func(t *testing.T) {
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `[]`))
		assert.NoError(t, err)

		result, err := cli.Country("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'countryCode' is required")
		assert.Nil(t, result)
	})

	t.Run("successful country request", func(t *testing.T) {
		cli, err := client.New("test-api-key", newMockClient(http.StatusOK, `[{"code": "GB-ENG", "name": "England"}]`))
		assert.NoError(t, err)

		result, err := cli.Country("GB-ENG")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "GB-ENG", result.Code)
	})
}

func TestLeague(t *testing.T) {
	responseBody := `[
		{
			"id": 2486,
			"name": "UEFA Champions League",
			"seasons": [{"season": 2025}, {"season": 2026}],
			"country": {"code": "World", "name": "World"}
		}
	]`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.League(2486)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "UEFA Champions League", result.Name)
	assert.Len(t, result.Seasons, 2)
	assert.Equal(t, 2025, result.Seasons[0].Season)
}

func TestLineups(t *testing.T) {
	responseBody := `{
		"homeTeam": {
			"id": 36526,
			"name": "Arsenal",
			"formation": "4-3-3",
			"initialLineup": [{"id": 235214, "name": "Bukayo Saka"}],
			"substitutes": [{"id": 12345, "name": "Leandro Trossard"}]
		},
		"awayTeam": {"id": 43334, "name": "Manchester City", "formation": "4-2-3-1", "initialLineup": [], "substitutes": []}
	}`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.Lineups(1281156605)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Arsenal", result.HomeTeam.Name)
	assert.Equal(t, "4-3-3", result.HomeTeam.Formation)
	assert.Equal(t, "Bukayo Saka", result.HomeTeam.InitialLineup[0].Name)
}

func TestEvents(t *testing.T) {
	responseBody := `[
		{"team": {"id": 1989571, "name": "Ñublense"}, "time": "23", "type": "Goal", "assist": "M. Plaza",
		 "player": "S. Valencia", "playerId": 91557473, "substituted": null, "assistingPlayerId": 49908061},
		{"team": {"id": 1973402, "name": "Palestino"}, "time": "46", "type": "Substitution", "assist": null,
		 "player": "B. Carrasco", "playerId": 5795671, "substituted": "M. Eggel", "assistingPlayerId": 47340916}
	]`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.Events(1281156605)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	events := *result
	assert.Len(t, events, 2)
	assert.Equal(t, "Goal", events[0].Type)
	assert.Equal(t, "M. Eggel", events[1].Substituted)
}

func TestPlayerStatistics(t *testing.T) {
	responseBody := `[
		{
			"id": 235214,
			"name": "Bukayo Saka",
			"fullName": "Bukayo Ayoyinka Temidayo Saka",
			"perCompetition": [
				{"club": "Arsenal FC", "type": "national league", "goals": 7, "league": "Premier League",
				 "season": "25/26", "assists": 7, "gamesPlayed": 31, "minutesPlayed": 2226}
			]
		}
	]`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.PlayerStatistics(235214)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Bukayo Saka", result.Name)
	assert.Len(t, result.PerCompetition, 1)
	assert.Equal(t, "25/26", result.PerCompetition[0].Season)
	assert.Equal(t, 7, result.PerCompetition[0].Goals)
}

func TestOdds(t *testing.T) {
	responseBody := `{
		"data": [
			{
				"match": {"id": 1281156605},
				"bookmaker": {"id": 1, "name": "Stake.com"},
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
		"plan": {"tier": "PRO", "message": "All data available with current plan."},
		"pagination": {"totalCount": 1, "offset": 0, "limit": 1}
	}`
	cli, err := client.New("test-api-key", newMockClient(http.StatusOK, responseBody))
	assert.NoError(t, err)

	result, err := cli.Odds(map[string]any{"matchId": 1281156605})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Stake.com", result.Data[0].Bookmaker.Name)
	assert.Equal(t, "1.85", result.Data[0].Markets[0].Odds[0].Value)
}
