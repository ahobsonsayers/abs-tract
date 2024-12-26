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
	"sync"

	"github.com/ahobsonsayers/abs-tract/utils"
	"github.com/orsinium-labs/enum"
	"github.com/samber/lo"
)

const (
	DefaultGoodreadsUrl = "https://www.goodreads.com"
	DefaultAPIKey       = "ckvsiSDsuqh7omh74ZZ6Q" // Read only API key kindly provided by LazyLibrarian
)

var (
	defaultGoodreadsUrl = lo.Must(url.Parse(DefaultGoodreadsUrl))

	DefaultClient = &Client{
		client:       http.DefaultClient,
		goodreadsUrl: utils.CloneURL(defaultGoodreadsUrl),
		apiKey:       DefaultAPIKey,
	}

	bookSearchTypeEnum   = enum.NewBuilder[string, BookSearchType]()
	BookSearchTypeTitle  = bookSearchTypeEnum.Add(BookSearchType{"title"})
	BookSearchTypeAuthor = bookSearchTypeEnum.Add(BookSearchType{"author"})
	BookSearchTypes      = bookSearchTypeEnum.Enum()
)

type BookSearchType enum.Member[string]

type Client struct {
	client       *http.Client
	goodreadsUrl *url.URL
	apiKey       string
}

// URL returns a clone of of the amazon url used by the client
func (c *Client) URL() *url.URL { return utils.CloneURL(c.goodreadsUrl) }

func (c *Client) get(
	ctx context.Context,
	path string,
	parameters map[string]string,
	target any,
) error {
	queryParams := url.Values{}
	queryParams.Add("key", c.apiKey)
	for key, value := range parameters {
		// Do some minor sanitising of value by removing characters known to be an issue
		sanitisedValue := regexp.MustCompile(`[\{\}]+`).ReplaceAllString(value, "")
		queryParams.Add(key, sanitisedValue)
	}

	requestUrl := c.URL()
	requestUrl = requestUrl.JoinPath(path)
	requestUrl.RawQuery = queryParams.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer response.Body.Close()

	httpError := utils.HTTPResponseError(response)
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
	err := c.get(ctx, "book/show.xml", queryParams, &result)
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

	// Only return books whose work have a title
	validBooks := make([]Book, 0, len(books))
	for _, book := range books {
		if book.Work.Title != "" {
			validBooks = append(validBooks, book)
		}
	}

	return validBooks, nil
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
	err := c.get(ctx, "book/title.xml", queryParams, &result)
	if err != nil {
		return Book{}, err
	}

	return result.Work, nil
}

// SearchBooks search for a book by its title and optionally an author (which can give better results)
// https://www.goodreads.com/api/index#search.books
func (c *Client) SearchBooks(ctx context.Context, bookTitle string, bookAuthor *string) ([]Book, error) {
	var bookOverviews []BookOverview
	var err error
	if bookAuthor == nil || *bookAuthor == "" {
		// If author is not set, search for books by title
		bookOverviews, err = c.searchBooksByTitle(ctx, bookTitle)
	} else {
		bookOverviews, err = c.searchBooksByTitleAndAuthor(ctx, bookTitle, *bookAuthor)
	}
	if err != nil {
		return nil, err
	}

	// Get book details using their ids
	bookDetails, err := c.GetBooksByIds(ctx, BookIds(bookOverviews))
	if err != nil {
		return nil, err
	}

	return bookDetails, nil
}
