package goodreads

import (
	"net/http"
	"strings"
)

// NewClient creates a new goodreads client.
// If client is nil, the default http client will be used.
// If api url is nil or uset, the default goodreads api url will be used
func NewClient(client *http.Client, apiURL *string, apiKey *string) *GoodreadsClient {
	if client == nil {
		client = http.DefaultClient
	}

	apiUrlString := DefaultAPIRootUrl
	if apiURL != nil && *apiURL != "" {
		apiUrlString = strings.Trim(*apiURL, "/")
	}

	apiKeyString := DefaultAPIKey
	if apiKey != nil && *apiKey != "" {
		apiKeyString = strings.TrimSpace(*apiKey)
	}

	return &GoodreadsClient{
		client:     client,
		apiRootUrl: apiUrlString,
		apiKey:     apiKeyString,
	}
}
