package kindle

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ahobsonsayers/abs-tract/utils"
	"github.com/samber/lo"
	"golang.org/x/net/html"
)

const DefaultAmazonURL = "https://www.amazon.com"

var (
	defaultAmazonURL = lo.Must(url.Parse(DefaultAmazonURL))

	DefaultClient = &Client{
		client:    http.DefaultClient,
		amazonUrl: utils.CloneURL(defaultAmazonURL),
	}
)

type Client struct {
	client    *http.Client
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

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "") // Amazon blocks some user agents

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, nil
	}
	defer response.Body.Close()

	httpError := utils.HTTPResponseError(response)
	if httpError != nil {
		return nil, httpError
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
