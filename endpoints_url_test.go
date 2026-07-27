package client_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	client "github.com/0ffsideCompass/highlightly-go-client"
)

// TestEndpointURLs verifies that every method requests the expected URL.
func TestEndpointURLs(t *testing.T) {
	tests := []struct {
		name        string
		call        func(cli *client.Client) error
		expectedURL string
		// body is the mock response; single resources decode from an array,
		// list endpoints from their wrapper (or a bare array).
		body string
	}{
		{
			name: "Countries",
			call: func(cli *client.Client) error {
				_, err := cli.Countries(map[string]any{"name": "England"})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/countries?name=England",
			body:        `[]`,
		},
		{
			name: "Country",
			call: func(cli *client.Client) error {
				_, err := cli.Country("GB")
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/countries/GB",
			body:        `[{}]`,
		},
		{
			name: "Leagues",
			call: func(cli *client.Client) error {
				_, err := cli.Leagues(map[string]any{"season": 2023})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/leagues?season=2023",
		},
		{
			name: "League",
			call: func(cli *client.Client) error {
				_, err := cli.League(2486)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/leagues/2486",
			body:        `[{}]`,
		},
		{
			name: "Highlights",
			call: func(cli *client.Client) error {
				_, err := cli.Highlights(map[string]any{"matchId": 1028043})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/highlights?matchId=1028043",
		},
		{
			name: "Highlight",
			call: func(cli *client.Client) error {
				_, err := cli.Highlight(99)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/highlights/99",
			body:        `[{}]`,
		},
		{
			name: "HighlightGeoRestrictions",
			call: func(cli *client.Client) error {
				_, err := cli.HighlightGeoRestrictions(99)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/highlights/geo-restrictions/99",
		},
		{
			name: "Matches",
			call: func(cli *client.Client) error {
				_, err := cli.Matches(map[string]any{"date": "2023-08-12"})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/matches?date=2023-08-12",
		},
		{
			name: "Match",
			call: func(cli *client.Client) error {
				_, err := cli.Match(1028043)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/matches/1028043",
			body:        `[{}]`,
		},
		{
			name: "LastFiveGames",
			call: func(cli *client.Client) error {
				_, err := cli.LastFiveGames(map[string]any{"teamId": 3057})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/last-five-games?teamId=3057",
			body:        `[]`,
		},
		{
			name: "HeadToHead",
			call: func(cli *client.Client) error {
				_, err := cli.HeadToHead(map[string]any{"teamIdOne": 3057, "teamIdTwo": 3068})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/head-2-head?teamIdOne=3057&teamIdTwo=3068",
			body:        `[]`,
		},
		{
			name: "Lineups",
			call: func(cli *client.Client) error {
				_, err := cli.Lineups(1028043)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/lineups/1028043",
		},
		{
			name: "Events",
			call: func(cli *client.Client) error {
				_, err := cli.Events(1028043)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/events/1028043",
			body:        `[]`,
		},
		{
			name: "Bookmakers",
			call: func(cli *client.Client) error {
				_, err := cli.Bookmakers(map[string]any{"name": "bet365"})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/bookmakers?name=bet365",
		},
		{
			name: "Bookmaker",
			call: func(cli *client.Client) error {
				_, err := cli.Bookmaker(1)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/bookmakers/1",
			body:        `[{}]`,
		},
		{
			name: "Odds",
			call: func(cli *client.Client) error {
				_, err := cli.Odds(map[string]any{"matchId": 1028043})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/odds?matchId=1028043",
		},
		{
			name: "Standings",
			call: func(cli *client.Client) error {
				_, err := cli.Standings(map[string]any{"leagueId": 2486, "season": 2023})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/standings?leagueId=2486&season=2023",
		},
		{
			name: "TeamStatistics",
			call: func(cli *client.Client) error {
				_, err := cli.TeamStatistics(3057, map[string]any{"fromDate": "2023-08-12"})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/teams/statistics/3057?fromDate=2023-08-12",
			body:        `[]`,
		},
		{
			name: "MatchStatistics",
			call: func(cli *client.Client) error {
				_, err := cli.MatchStatistics(1028043)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/statistics/1028043",
			body:        `[]`,
		},
		{
			name: "Players",
			call: func(cli *client.Client) error {
				_, err := cli.Players(map[string]any{"name": "Saka"})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/players?name=Saka",
		},
		{
			name: "Player",
			call: func(cli *client.Client) error {
				_, err := cli.Player(5001)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/players/5001",
			body:        `[{}]`,
		},
		{
			name: "PlayerStatistics",
			call: func(cli *client.Client) error {
				_, err := cli.PlayerStatistics(5001)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/players/5001/statistics",
			body:        `[{}]`,
		},
		{
			name: "Teams",
			call: func(cli *client.Client) error {
				_, err := cli.Teams(map[string]any{"name": "Arsenal"})
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/teams?name=Arsenal",
		},
		{
			name: "Team",
			call: func(cli *client.Client) error {
				_, err := cli.Team(3057)
				return err
			},
			expectedURL: "https://sports.highlightly.net/football/teams/3057",
			body:        `[{}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := tt.body
			if body == "" {
				body = `{}`
			}
			mock := newMockClient(http.StatusOK, body)
			cli, err := client.New("test-api-key", mock)
			assert.NoError(t, err)

			assert.NoError(t, tt.call(cli))
			assert.NotNil(t, mock.LastRequest)
			assert.Equal(t, tt.expectedURL, mock.LastRequest.URL.String())
		})
	}
}
