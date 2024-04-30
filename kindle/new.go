package kindle

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ahobsonsayers/abs-tract/utils"
	"github.com/samber/lo"
)

var countryAmazonURLs = map[string]*url.URL{
	"au": lo.Must(url.Parse("https://www.amazon.com.au")),
	"ca": lo.Must(url.Parse("https://www.amazon.ca")),
	"de": lo.Must(url.Parse("https://www.amazon.de")),
	"es": lo.Must(url.Parse("https://www.amazon.es")),
	"fr": lo.Must(url.Parse("https://www.amazon.fr")),
	"in": lo.Must(url.Parse("https://www.amazon.co.in")),
	"it": lo.Must(url.Parse("https://www.amazon.it")),
	"jp": lo.Must(url.Parse("https://www.amazon.co.jp")),
	"uk": lo.Must(url.Parse("https://www.amazon.co.uk")),
	"us": defaultAmazonURL,
}

// If client is nil, the default http client will be used.
// If country code is nil or unset, amazon.com will be used as the url.
// Will return an error if the country code is invalid.
func NewClient(client *http.Client, countryCode *string) (*Client, error) {
	if client == nil {
		client = http.DefaultClient
	}

	amazonUrl := defaultAmazonURL
	if countryCode != nil && *countryCode != "" {
		countryAmazonUrl, ok := countryAmazonURLs[strings.Trim(*countryCode, "/")]
		if !ok {
			return nil, fmt.Errorf("invalid country code: %s", *countryCode)
		}
		amazonUrl = countryAmazonUrl
	}

	return &Client{
		client:    client,
		amazonUrl: utils.CloneURL(amazonUrl),
	}, nil
}
