package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	// domain is the default base URL. The All Sports host namespaces every
	// resource under a sport segment; per-sport hosts (for example
	// https://soccer.highlightly.net/) serve the same resources without it,
	// so the sport prefix lives in the domain rather than in endpoint paths.
	domain = "https://sports.highlightly.net/football/"
)

const (
	authKey = "x-rapidapi-key"
	hostKey = "x-rapidapi-host"
	hostVal = "sport-highlights-api.p.rapidapi.com"
)

// HttpClient is an interface that abstracts the http.Client's Do method,
// allowing for easier testing and customization.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents a client for the Highlightly sports API.
// It holds the API key, the base domain URL, and the underlying HTTP client used to execute requests.
type Client struct {
	key    string
	Domain string
	client HttpClient
}

// New creates a new Client instance for the Highlightly API using the default
// football domain. It returns an error if the API key is missing or the
// provided HTTP client is nil.
func New(key string, client HttpClient) (*Client, error) {
	return NewWithDomain(key, domain, client)
}

// NewWithDomain creates a new Client instance with a custom domain, such as a
// per-sport host (https://soccer.highlightly.net/) or the RapidAPI gateway.
// It returns an error if the API key or domain is missing, or the provided
// HTTP client is nil.
func NewWithDomain(key, domain string, client HttpClient) (*Client, error) {
	if key == "" {
		return nil, errors.New("missing key")
	}
	if domain == "" {
		return nil, errors.New("missing domain")
	}
	if client == nil {
		return nil, errors.New("missing http client")
	}
	// Endpoint paths are appended directly to the domain, so it must end with a slash.
	if !strings.HasSuffix(domain, "/") {
		domain += "/"
	}
	return &Client{
		key:    key,
		Domain: domain,
		client: client,
	}, nil
}

func (c *Client) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add(authKey, c.key)
	// Highlightly's own hosts ignore this header; the RapidAPI gateway
	// requires it, so it is safe to send in both modes.
	req.Header.Add(hostKey, hostVal)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Check for non-2xx status codes
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		var body []byte
		if b, readErr := io.ReadAll(res.Body); readErr == nil {
			body = b
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", res.StatusCode, string(body))
	}

	return io.ReadAll(res.Body)
}

// paramString formats a query parameter value. JSON-decoded numbers arrive as
// float64 and must render as plain integers when whole (1581037), never in
// scientific notation (1.581037e+06), which the API rejects.
func paramString(value any) string {
	switch v := value.(type) {
	case float64:
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// buildURL constructs the URL with dynamic parameters using net/url
func (c *Client) buildURL(endpoint string, params map[string]any) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		return endpoint
	}

	values := url.Values{}
	for key, value := range params {
		values.Add(key, paramString(value))
	}

	u.RawQuery = values.Encode()
	return u.String()
}
