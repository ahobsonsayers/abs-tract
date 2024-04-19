package goodreads

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"

	"github.com/ahobsonsayers/abs-goodreads/utils"
)

const (
	DefaultAPIRootUrl = "https://www.goodreads.com"
	DefaultAPIKey     = "ckvsiSDsuqh7omh74ZZ6Q" // Read only API key kindly provided by LazyLibrarian
)

var DefaultClient = &Client{
	client:     http.DefaultClient,
	apiRootUrl: DefaultAPIRootUrl,
	apiKey:     DefaultAPIKey,
}

type Client struct {
	client     *http.Client
	apiRootUrl string
	apiKey     string
}

func (c *Client) Get(
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
		// Do some minor sanitising of value by removing characters known to be an issue
		value = regexp.MustCompile(`[\{\}]+`).ReplaceAllString(value, "")
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

// GetBookById gets a book by its id.
// https://www.goodreads.com/api/index#book.show
func (c *Client) GetBookById(ctx context.Context, bookId string) (Book, error) {
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

func (c *Client) GetBooksByIds(ctx context.Context, bookIds []string) ([]Book, error) {
	books := make([]Book, len(bookIds))
	var errs error

	var wg sync.WaitGroup
	var booksMutex sync.Mutex
	var errsMutex sync.Mutex

	for idx, bookId := range bookIds {
		wg.Add(1)

		go func(bookId string, idx int) {
			defer wg.Done()

			book, err := c.GetBookById(ctx, bookId)
			if err != nil {
				errsMutex.Lock()
				errs = errors.Join(errs, err)
				errsMutex.Unlock()
				return
			}

			booksMutex.Lock()
			books[idx] = book
			booksMutex.Unlock()
		}(bookId, idx)
	}

	wg.Wait()

	if errs != nil {
		return nil, errs
	}

	return books, nil
}

// GetBookByTitle gets a book by its title and optionally an author (which can give a better match)
// https://www.goodreads.com/api/index#book.title
func (c *Client) GetBookByTitle(ctx context.Context, bookTitle string, bookAuthor *string) (Book, error) {
	queryParams := map[string]string{"title": bookTitle}
	if bookAuthor != nil && *bookAuthor != "" {
		queryParams["author"] = *bookAuthor
	}

	var result struct {
		Work Book `xml:"book"`
	}
	err := c.Get(ctx, "book/title.xml", queryParams, &result)
	if err != nil {
		return Book{}, err
	}

	return result.Work, nil
}

// SearchBooks search for a book by its title and optionally an author (which can give better results)
// https://www.goodreads.com/api/index#search.books
func (c *Client) SearchBooks(ctx context.Context, bookTitle string, bookAuthor *string) ([]Book, error) {
	// Search for books via title, getting their ids.
	bookIds, err := c.searchBookIdsByTitle(ctx, bookTitle)
	if err != nil {
		return nil, err
	}

	// If author is set, also search for books via author, getting their ids.
	// Only keeps book ids that appear in both title and author searches
	if bookAuthor != nil && *bookAuthor != "" {
		authorBookIds, err := c.searchBookIdsByAuthor(ctx, *bookAuthor)
		if err != nil {
			return nil, err
		}

		// Get common book ids. If there are no common book
		// ids, just use the book ids form the title search.
		// Some result are better than none!
		commonBookIds := utils.Intersection(bookIds, authorBookIds)
		if len(commonBookIds) != 0 {
			bookIds = commonBookIds
		}

	}

	// Get book details using their ids
	books, err := c.GetBooksByIds(ctx, bookIds)
	if err != nil {
		return nil, err
	}

	return books, nil
}

type bookIdUnmarshaller struct {
	BookIds []string `xml:"search>results>work>best_book>id"`
}

func (c *Client) searchBookIdsByTitle(ctx context.Context, title string) ([]string, error) {
	queryParams := map[string]string{
		"q":             title,
		"search[field]": "title",
	}
	var unmarshaller bookIdUnmarshaller
	err := c.Get(ctx, "search/index.xml", queryParams, &unmarshaller)
	if err != nil {
		return nil, err
	}

	return unmarshaller.BookIds, nil
}

func (c *Client) searchBookIdsByAuthor(ctx context.Context, author string) ([]string, error) {
	var bookIds []string
	var errs error

	var wg sync.WaitGroup
	var bookIdsMutex sync.Mutex
	var errsMutex sync.Mutex

	// Get 10 pages of author books. This should be enough
	for pageNumber := 1; pageNumber <= 10; pageNumber++ {
		wg.Add(1)

		go func(page int) {
			defer wg.Done()

			queryParams := map[string]string{
				"q":             author,
				"page":          strconv.Itoa(page),
				"search[field]": "author",
			}
			var unmarshaller bookIdUnmarshaller
			err := c.Get(ctx, "search/index.xml", queryParams, &unmarshaller)
			if err != nil {
				errsMutex.Lock()
				errs = errors.Join(errs, err)
				errsMutex.Unlock()
				return
			}

			bookIdsMutex.Lock()
			bookIds = append(bookIds, unmarshaller.BookIds...)
			bookIdsMutex.Unlock()
		}(pageNumber)
	}
	wg.Wait()

	return bookIds, errs
}
