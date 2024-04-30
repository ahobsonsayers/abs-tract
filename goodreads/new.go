package goodreads

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ahobsonsayers/abs-tract/utils"
)

// NewClient creates a new goodreads client.
// If client is nil, the default http client will be used.
// If goodreads url is nil or unset, the default goodreads url will be used.
// Will return an error if the goodreads url is invalid.
func NewClient(client *http.Client, goodreadsUrl, apiKey *string) (*Client, error) {
	if client == nil {
		client = http.DefaultClient
	}

	goodreadsUrlStruct := defaultGoodreadsUrl
	if goodreadsUrl != nil && *goodreadsUrl != "" {
		parsedGoodreadsUrl, err := url.Parse(strings.Trim(*goodreadsUrl, "/"))
		if err != nil {
			return nil, fmt.Errorf("invalid goodreads url: %w", err)
		}
		goodreadsUrlStruct = parsedGoodreadsUrl
	}

	apiKeyString := DefaultAPIKey
	if apiKey != nil && *apiKey != "" {
		apiKeyString = strings.TrimSpace(*apiKey)
	}

	return &Client{
		client:       client,
		goodreadsUrl: utils.CloneURL(goodreadsUrlStruct),
		apiKey:       apiKeyString,
	}, nil
}
