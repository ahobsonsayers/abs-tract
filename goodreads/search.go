package goodreads

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/samber/lo"
)

// SearchBooks search for a book by its title and optionally an author (which can give better results)
// https://www.goodreads.com/api/index#search.books
func (c *Client) SearchBooks(ctx context.Context, title string, author *string) ([]Book, error) {
	if author == nil || *author == "" {
		// If author is not set, search for books by title
		return c.searchBooksByTitle(ctx, title)
	}

	return c.searchBooksByTitleAndAuthor(ctx, title, *author)
}

func (c *Client) searchBooksByTitle(ctx context.Context, title string) ([]Book, error) {
	bookOverviews, err := c.searchBooksManyPages(ctx, searchBooksManyPagesInput{
		Query:      title,
		SearchType: BookSearchTypeTitle,
		Page:       1,
		NumPages:   5,
	})
	if err != nil {
		return nil, err
	}

	return c.GetBooksByIds(ctx, BookIds(bookOverviews))
}

func (c *Client) searchBooksByTitleAndAuthor(
	ctx context.Context,
	title string,
	author string,
) ([]Book, error) {
	// If searching by title and author, search by author ONLY first to get their books.
	// We will then search the books for title.
	// We do NOT search goodreads by title AND author together using the 'all' search type
	// as goodreads returns awful results, including sometimes none at all.
	bookOverviews, err := c.searchBooksManyPages(ctx, searchBooksManyPagesInput{
		Query:      author,
		SearchType: BookSearchTypeAuthor,
		Page:       1,
		NumPages:   15,
	})
	if err != nil {
		return nil, err
	}

	books, err := c.GetBooksByIds(ctx, BookIds(bookOverviews))
	if err != nil {
		return nil, err
	}

	sortBookByTitleSimilarity(books, title)

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
