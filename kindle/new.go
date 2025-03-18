package kindle

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ahobsonsayers/abs-tract/utils"
	"github.com/imroc/req/v3"
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

// Creates a new kindle client.
// If country code is nil or unset, amazon.com will be used as the url.
// Will return an error if the country code is invalid.
func NewClient(countryCode *string) (*Client, error) {
	client := req.C().ImpersonateFirefox()

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
