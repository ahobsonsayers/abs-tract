package goodreads

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	DefaultAPIRootUrl = "https://www.goodreads.com"
	DefaultAPIKey     = "ckvsiSDsuqh7omh74ZZ6Q" // Read only API key kindly provided by LazyLibrarian
)

var DefaultGoodreadsClient = &GoodreadsClient{
	client:     http.DefaultClient,
	apiRootUrl: DefaultAPIRootUrl,
	apiKey:     DefaultAPIKey,
}

type GoodreadsClient struct {
	client     *http.Client
	apiRootUrl string
	apiKey     string
}

func (c *GoodreadsClient) Get(
	ctx context.Context,
	apiPath string,
	queryParams map[string]string,
	target any,
) error {
	// Construct api url
	apiUrl, err := url.Parse(c.apiRootUrl)
	if err != nil {
		return fmt.Errorf("failed to parse api root url: %w", err)
	}
	apiUrl = apiUrl.JoinPath(apiPath)

	apiUrlValues := make(url.Values, len(queryParams))
	apiUrlValues.Add("key", c.apiKey)
	for key, value := range queryParams {
		apiUrlValues.Add(key, value)
	}
	apiUrl.RawQuery = url.Values(apiUrlValues).Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer response.Body.Close()

	httpError := HTTPResponseError(response)
	if httpError != nil {
		return httpError
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	err = xml.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}

// GetBook gets a book by its id.
// https://www.goodreads.com/api/index#book.show
func (c *GoodreadsClient) GetBook(ctx context.Context, bookId string) (Book, error) {
	queryParams := map[string]string{"id": bookId}

	var result struct {
		Book Book `xml:"book"`
	}
	err := c.Get(ctx, "book/show.xml", queryParams, &result)
	if err != nil {
		return Book{}, err
	}

	return result.Book, nil
}

// SearchBook search for a book by its title. An author can be specified to give better results.
// https://www.goodreads.com/api/index#book.title
func (c *GoodreadsClient) SearchBook(ctx context.Context, bookTitle string, bookAuthor *string) ([]Book, error) {
	queryParams := map[string]string{"title": bookTitle}
	if bookAuthor != nil && *bookAuthor != "" {
		queryParams["author"] = *bookAuthor
	}

	var result struct {
		Works []Book `xml:"book"`
	}
	err := c.Get(ctx, "book/title.xml", queryParams, &result)
	if err != nil {
		return nil, err
	}

	return result.Works, nil
}
