package goodreads

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/orsinium-labs/enum"
	"github.com/samber/lo"
)

var (
	bookSearchTypeEnum = enum.NewBuilder[string, BookSearchType]()

	BookSearchTypeTitle  = bookSearchTypeEnum.Add(BookSearchType{"title"})
	BookSearchTypeAuthor = bookSearchTypeEnum.Add(BookSearchType{"author"})
	BookSearchTypeAll    = bookSearchTypeEnum.Add(BookSearchType{"all"})

	BookSearchTypes = bookSearchTypeEnum.Enum()
)

// SearchBooks search for a book by its title and optionally an author (which can give better ordered results).
// Returns the first 10 pages of books.
// See: https://www.goodreads.com/api/index#search.books
func (c *Client) SearchBooks(ctx context.Context, title string, author *string) ([]Book, error) {
	// Normalise title and author to make searching more consistent
	normalisedTitle := normaliseString(title)
	normalisedAuthor := normaliseString(lo.FromPtr(author))

	switch {
	case normalisedTitle != "" && normalisedAuthor != "":
		// If both title and author are set
		return c.searchBooksByTitleAndAuthor(ctx, title, normalisedAuthor)

	case normalisedTitle != "":
		// If only title is set
		return c.searchBooksByTitle(ctx, normalisedTitle)

	case normalisedAuthor != "":
		// If only author is set
		return c.searchBooksByAuthor(ctx, normalisedAuthor)
	}

	return nil, nil
}

func (c *Client) searchBooksByTitle(ctx context.Context, title string) ([]Book, error) {
	bookOverviews, err := c.searchBooksManyPages(ctx, searchBooksManyPagesInput{
		Query:      title,
		SearchType: BookSearchTypeTitle,
		Page:       1,
		NumPages:   10,
	})
	if err != nil {
		return nil, err
	}

	return c.GetBooksByIds(ctx, BookIds(bookOverviews))
}

func (c *Client) searchBooksByAuthor(ctx context.Context, author string) ([]Book, error) {
	bookOverviews, err := c.searchBooksManyPages(ctx, searchBooksManyPagesInput{
		Query:      author,
		SearchType: BookSearchTypeAuthor,
		Page:       1,
		NumPages:   10,
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
	// If searching by title and author, search by title ONLY first.
	// We will then sort books but author similarity.
	// We do NOT search goodreads by title AND author together as goodreads
	// returns awful results, including sometimes none at all.
	books, err := c.searchBooksByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	sortBookByAuthorSimilarity(books, author)

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
		// Default to all search
		input.SearchType = BookSearchTypeAll
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

	return unmarshaller.Books, nil
}
