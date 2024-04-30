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
	"strings"
	"sync"

	"github.com/ahobsonsayers/abs-tract/utils"
	"github.com/lithammer/fuzzysearch/fuzzy"
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
		bookOverviews, err = c.searchBooksManyPages(ctx, searchBooksManyPagesInput{
			Query:      bookTitle,
			SearchType: BookSearchTypeTitle,
			Page:       1,
			NumPages:   5,
		})
		if err != nil {
			return nil, err
		}

	} else {
		// If author is set, search for books by author ONLY.
		// We will then search the authors books for title.
		// We do NOT search goodreads by title AND author together using the 'all' search type
		// as goodreads returns awful results, including sometimes none at all.
		authorBookOverviews, err := c.searchBooksManyPages(ctx, searchBooksManyPagesInput{
			Query:      *bookAuthor,
			SearchType: BookSearchTypeAuthor,
			Page:       1,
			NumPages:   15,
		})
		if err != nil {
			return nil, err
		}

		// In author books, find books that (fuzzily) match the title
		authorBookOverviewMatches := fuzzy.RankFindNormalizedFold(bookTitle, BookTitles(authorBookOverviews))

		// Use matched author books for the book overviews. User order returned by goodreads
		matchedBookOverviews := make([]BookOverview, 0, len(authorBookOverviewMatches))
		for _, match := range authorBookOverviewMatches {
			matchedBookOverviews = append(matchedBookOverviews, authorBookOverviews[match.OriginalIndex])
		}
		bookOverviews = matchedBookOverviews
	}

	// Get book details using their ids
	bookDetails, err := c.GetBooksByIds(ctx, BookIds(bookOverviews))
	if err != nil {
		return nil, err
	}

	return bookDetails, nil
}

type searchBooksSinglePageInput struct {
	Query      string
	SearchType BookSearchType
	Page       int
}

// searchBooksPage searches for books and returns the requested page of results.
// If Query is not set, no search will be performed and no results will be returned.
// If search type is unset or invalid, search will fallback to a title search
// If page is < 1, first page of results is returned.
func (c *Client) searchBooksSinglePage(ctx context.Context, input searchBooksSinglePageInput) ([]BookOverview, error) {
	input.Query = strings.TrimSpace(input.Query)
	if input.Query == "" {
		return nil, nil
	}
	if !BookSearchTypes.Contains(input.SearchType) {
		// Default to title search
		input.SearchType = BookSearchTypeTitle
	}
	if input.Page < 1 {
		input.Page = 1
	}

	queryParams := map[string]string{
		"q":             input.Query,
		"search[field]": input.SearchType.Value,
		"page":          strconv.Itoa(input.Page),
	}
	var unmarshaller struct {
		Books []BookOverview `xml:"search>results>work>best_book"`
	}
	err := c.get(ctx, "search/index.xml", queryParams, &unmarshaller)
	if err != nil {
		return nil, err
	}

	// Sanitise the books
	books := make([]BookOverview, 0, len(unmarshaller.Books))
	for _, book := range unmarshaller.Books {
		book.Sanitise()
		books = append(books, book)
	}

	return books, nil
}

type searchBooksManyPagesInput struct {
	Query      string
	SearchType BookSearchType
	Page       int
	NumPages   int
}

// searchBooksManyPages searches for books and returns (flattened) results from the request number of pages.
// Arguments are the same as searchBooks except for NumPages. If NumPages < 1, a single page of results
// will be returned
func (c *Client) searchBooksManyPages(ctx context.Context, input searchBooksManyPagesInput) ([]BookOverview, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.NumPages < 1 {
		input.NumPages = 1
	}

	bookPages := make([][]BookOverview, input.NumPages)
	var errs error

	var wg sync.WaitGroup
	var bookIdsMutex sync.Mutex
	var errsMutex sync.Mutex

	// Get 5 pages of books. This should be enough
	for idx := 0; idx < input.NumPages; idx++ {
		wg.Add(1)

		go func(idx int) {
			defer wg.Done()

			pageBooks, err := c.searchBooksSinglePage(
				ctx, searchBooksSinglePageInput{
					Query:      input.Query,
					SearchType: input.SearchType,
					Page:       input.Page + idx,
				},
			)
			if err != nil {
				errsMutex.Lock()
				errs = errors.Join(errs, err)
				errsMutex.Unlock()
				return
			}

			bookIdsMutex.Lock()
			bookPages[idx] = pageBooks
			bookIdsMutex.Unlock()
		}(idx)
	}
	wg.Wait()

	return lo.Flatten(bookPages), errs
}
