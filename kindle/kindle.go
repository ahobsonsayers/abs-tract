package kindle

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ahobsonsayers/abs-tract/utils"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"golang.org/x/net/html"
)

const DefaultAmazonURL = "https://www.amazon.com"

var (
	defaultAmazonURL = lo.Must(url.Parse(DefaultAmazonURL))

	DefaultClient = &Client{
		client:    req.C().ImpersonateFirefox(),
		amazonUrl: utils.CloneURL(defaultAmazonURL),
	}
)

type Client struct {
	client    *req.Client
	amazonUrl *url.URL
}

// URL returns a clone of of the amazon url used by the client
func (c *Client) URL() *url.URL { return utils.CloneURL(c.amazonUrl) }

func (c *Client) get(
	ctx context.Context,
	path string,
	parameters map[string]string,
) (*html.Node, error) {
	queryParams := url.Values{}
	for key, value := range parameters {
		queryParams.Add(key, value)
	}

	requestUrl := c.URL()
	requestUrl = requestUrl.JoinPath(path)
	requestUrl.RawQuery = queryParams.Encode()

	response, err := c.client.R().SetContext(ctx).Get(requestUrl.String())
	if err != nil {
		return nil, nil
	}
	if !response.IsSuccessState() {
		return nil, fmt.Errorf("%s: %s", response.GetStatus(), response.String())
	}

	htmlResponse, err := html.Parse(response.Body)
	if err != nil {
		return nil, nil
	}

	return htmlResponse, nil
}

func (c *Client) Search(ctx context.Context, title string, author *string) ([]Book, error) {
	parameters := map[string]string{
		"i": "digital-text",
		// "i": "stripbooks",
		"k": title,
	}
	if author != nil && *author != "" {
		parameters["inauthor"] = *author
	}

	htmlResponse, err := c.get(ctx, "s", parameters)
	if err != nil {
		return nil, err
	}

	return BooksFromHTML(htmlResponse)
}
