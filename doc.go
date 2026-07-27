// Package client provides a Go client for the Highlightly sports API,
// covering every documented endpoint: countries, leagues, teams, matches,
// lineups, events, standings, statistics, players, odds, bookmakers, and
// video highlights with geo restrictions.
//
// # Creating a client
//
// A client needs an API key and an HTTP client. Any type implementing the
// HttpClient interface (a single Do method) works, so a plain *http.Client is
// enough:
//
//	cli, err := client.New(apiKey, &http.Client{Timeout: 10 * time.Second})
//	if err != nil {
//		log.Fatal(err)
//	}
//
// By default requests go to the All Sports host's football namespace at
// https://sports.highlightly.net/football/, authenticated with the
// x-rapidapi-key header. Endpoint paths are appended directly to the domain,
// so NewWithDomain accepts a per-sport host (https://soccer.highlightly.net/)
// or the RapidAPI gateway's equivalent base URL instead; the client sends the
// x-rapidapi-host header as well, so both providers work without further
// configuration.
//
// # Calling endpoints
//
// List methods take a map of query parameters mirroring the API's own
// parameter names, and return a typed response model from the models package:
//
//	highlights, err := cli.Highlights(map[string]any{
//		"leagueName": "Premier League",
//		"date":       "2023-08-12",
//		"limit":      10,
//	})
//
// Single-resource methods take the resource ID directly:
//
//	match, err := cli.Match(1028043)
//
// Each method's documentation lists the parameters the endpoint accepts,
// including which ones are required. Required parameters are validated before
// the request is sent.
//
// # Error handling
//
// Methods return wrapped errors for failed requests (non-2xx responses,
// with the response body included in the message) and for undecodable
// responses.
//
// # Retries and rate limits
//
// The client ships without retry logic by design: inject a retrying HTTP
// client (such as hashicorp/go-retryablehttp's StandardClient) if you need
// it. Note that Highlightly enforces request quotas per plan — see the Plan
// field on list responses for your current usage — so avoid blindly retrying
// 429 responses within the same rate window.
package client
