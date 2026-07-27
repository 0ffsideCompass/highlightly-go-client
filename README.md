# highlightly-go-client

A Go client for the [Highlightly](https://highlightly.net) sports API: live scores, matches, standings, team and player statistics, head-to-head records, odds, and on-demand video highlight clips.

## Installation

```bash
go get github.com/0ffsideCompass/highlightly-go-client
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	client "github.com/0ffsideCompass/highlightly-go-client"
)

func main() {
	cli, err := client.New("your-api-key", &http.Client{Timeout: 10 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	// List methods take a map of query parameters mirroring the API's own names.
	matches, err := cli.Matches(map[string]any{
		"leagueName": "Premier League",
		"date":       "2023-08-12",
		"limit":      10,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range matches.Data {
		fmt.Printf("%s vs %s (%s)\n", m.HomeTeam.Name, m.AwayTeam.Name, m.State.Score.Current)
	}

	// Single-resource methods take the resource ID directly.
	highlight, err := cli.Highlight(7701)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(highlight.Title, highlight.URL)
}
```

### Hosts

By default the client talks to the All Sports host's football namespace, `https://sports.highlightly.net/football/`. Endpoint paths are appended directly to the domain, so you can point the client at a per-sport host or the RapidAPI gateway with `NewWithDomain`:

```go
cli, err := client.NewWithDomain("your-api-key", "https://soccer.highlightly.net/", httpClient)
```

Authentication uses the `x-rapidapi-key` header; the `x-rapidapi-host` header is sent as well so calls through RapidAPI work without further configuration.

## Endpoints

| Method | Endpoint |
|---|---|
| `Countries(params)` | `GET /countries` |
| `Country(code)` | `GET /countries/{countryCode}` |
| `Leagues(params)` | `GET /leagues` |
| `League(id)` | `GET /leagues/{id}` |
| `Teams(params)` | `GET /teams` |
| `Team(id)` | `GET /teams/{id}` |
| `Matches(params)` | `GET /matches` |
| `Match(id)` | `GET /matches/{id}` |
| `LastFiveGames(params)` | `GET /last-five-games` |
| `HeadToHead(params)` | `GET /head-2-head` |
| `Lineups(matchID)` | `GET /lineups/{matchId}` |
| `Events(matchID)` | `GET /events/{id}` |
| `Standings(params)` | `GET /standings` |
| `TeamStatistics(teamID, params)` | `GET /teams/statistics/{id}` |
| `MatchStatistics(matchID)` | `GET /statistics/{matchId}` |
| `Players(params)` | `GET /players` |
| `Player(id)` | `GET /players/{id}` |
| `PlayerStatistics(id)` | `GET /players/{id}/statistics` |
| `Highlights(params)` | `GET /highlights` |
| `Highlight(id)` | `GET /highlights/{id}` |
| `HighlightGeoRestrictions(id)` | `GET /highlights/geo-restrictions/{id}` |
| `Bookmakers(params)` | `GET /bookmakers` |
| `Bookmaker(id)` | `GET /bookmakers/{id}` |
| `Odds(params)` | `GET /odds` |

Each method's Go doc lists the parameters the endpoint accepts, including which ones are required. Required parameters are validated before the request is sent.

## Testing

```bash
go test ./...
```

## License

See [LICENSE](LICENSE).
