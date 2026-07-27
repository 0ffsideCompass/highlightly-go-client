package client_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	client "github.com/0ffsideCompass/highlightly-go-client"
)

// MockHTTPClient is a mock implementation of the HttpClient interface that
// records the last request it received.
type MockHTTPClient struct {
	Response    *http.Response
	Err         error
	LastRequest *http.Request
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.LastRequest = req
	return m.Response, m.Err
}

func newMockClient(statusCode int, body string) *MockHTTPClient {
	return &MockHTTPClient{
		Response: &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
		},
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		client      client.HttpClient
		expectError string
	}{
		{
			name:   "valid client",
			key:    "test-api-key",
			client: &MockHTTPClient{},
		},
		{
			name:        "missing key",
			key:         "",
			client:      &MockHTTPClient{},
			expectError: "missing key",
		},
		{
			name:        "missing http client",
			key:         "test-api-key",
			client:      nil,
			expectError: "missing http client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := client.New(tt.key, tt.client)
			if tt.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
				assert.Nil(t, cli)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cli)
				assert.Equal(t, "https://sports.highlightly.net/football/", cli.Domain)
			}
		})
	}
}

func TestNewWithDomain(t *testing.T) {
	t.Run("missing domain", func(t *testing.T) {
		cli, err := client.NewWithDomain("test-api-key", "", &MockHTTPClient{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing domain")
		assert.Nil(t, cli)
	})

	t.Run("appends trailing slash", func(t *testing.T) {
		cli, err := client.NewWithDomain("test-api-key", "https://soccer.highlightly.net", &MockHTTPClient{})
		assert.NoError(t, err)
		assert.Equal(t, "https://soccer.highlightly.net/", cli.Domain)
	})

	t.Run("keeps trailing slash", func(t *testing.T) {
		cli, err := client.NewWithDomain("test-api-key", "https://soccer.highlightly.net/", &MockHTTPClient{})
		assert.NoError(t, err)
		assert.Equal(t, "https://soccer.highlightly.net/", cli.Domain)
	})
}

func TestAuthHeaders(t *testing.T) {
	mock := newMockClient(http.StatusOK, `{"data": []}`)
	cli, err := client.New("test-api-key", mock)
	assert.NoError(t, err)

	_, err = cli.Teams(nil)
	assert.NoError(t, err)
	assert.NotNil(t, mock.LastRequest)
	assert.Equal(t, "test-api-key", mock.LastRequest.Header.Get("x-rapidapi-key"))
	assert.Equal(t, "sport-highlights-api.p.rapidapi.com", mock.LastRequest.Header.Get("x-rapidapi-host"))
}

func TestHTTPClientError(t *testing.T) {
	mock := &MockHTTPClient{Err: errors.New("connection refused")}
	cli, err := client.New("test-api-key", mock)
	assert.NoError(t, err)

	result, err := cli.Teams(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
	assert.Nil(t, result)
}
