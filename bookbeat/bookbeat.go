// Package bookbeat provides functionality to search and retrieve book metadata from the BookBeat service.
package bookbeat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
)

// API endpoints for BookBeat services
const (
	SearchBaseURL = "https://www.bookbeat.com/api/next/search"
)

// breakTagRegex matches HTML break tags in book descriptions
var breakTagRegex = regexp.MustCompile(`<br\s*/?>`)

// languageMap converts language codes to BookBeat language identifiers
var languageMap = map[string]string{
	"en": "english",
	"de": "german",
	"ar": "arabic",
	"eu": "basque",
	"ca": "catalan",
	"cs": "czech",
	"da": "danish",
	"nl": "dutch",
	"et": "estonian",
	"fi": "finnish",
	"fr": "french",
	"hu": "hungarian",
	"it": "italian",
	"nb": "norwegian",
	"nn": "norwegiannynorsk",
	"pl": "polish",
	"pt": "portuguese",
	"ru": "russian",
	"es": "spanish",
	"sv": "swedish",
	"tr": "turkish",
}

// countryMap converts country codes to BookBeat market identifiers
var countryMap = map[string]int{
	"gr": 30,
	"nl": 31,
	"be": 32,
	"fr": 33,
	"es": 34,
	"hu": 36,
	"it": 39,
	"ro": 40,
	"ch": 41,
	"at": 43,
	"uk": 44,
	"dk": 45,
	"se": 46,
	"no": 47,
	"pl": 48,
	"de": 49,
	"pt": 351,
	"lu": 352,
	"ie": 353,
	"mt": 356,
	"cy": 357,
	"fi": 358,
	"bg": 359,
	"lt": 370,
	"lv": 371,
	"ee": 372,
	"hr": 385,
	"si": 386,
	"cz": 420,
	"sk": 421,
}

// Bookbeat represents a client for interacting with the BookBeat API
type Bookbeat struct {
	searchBaseURL string
	httpClient    *req.Client
	market        string
	formats       []string
	languages     []string
}

// Market returns set market
func (b *Bookbeat) Market() string {
	return b.market
}

// Languages returns set languages
func (b *Bookbeat) Languages() []string {
	return b.languages
}

// Formats returns set formats
func (b *Bookbeat) Formats() []string {
	return b.formats
}

// SearchBooks performs a search query against the BookBeat API and returns raw search results
func (b *Bookbeat) SearchBooks(ctx context.Context, query string, author *string) (SearchResponse, error) {
	params := url.Values{}
	params.Add("query", query)
	if author != nil {
		params.Add("author", *author)
	}
	params.Add("page", "1")
	params.Add("limit", "20")
	params.Add("sortby", "relevance")
	params.Add("sortorder", "desc")
	params.Add("includeerotic", "false")
	params.Add("market", strconv.Itoa(countryMap[b.market]))

	// Add format filters
	for _, format := range b.formats {
		params.Add("format", format)
	}

	// Add language filters
	for _, language := range b.languages {
		params.Add("language", language)
	}

	searchURL, _ := url.Parse(b.searchBaseURL)
	searchURL.RawQuery = params.Encode()

	response, err := b.httpClient.R().SetContext(ctx).Get(searchURL.String())
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	if !response.IsSuccessState() {
		return SearchResponse{},
			fmt.Errorf("search request failed with status %d: %s",
				response.StatusCode,
				response.Body)
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResp); err != nil {
		return SearchResponse{}, fmt.Errorf("failed to decode search response: %w", err)
	}
	searchResp.QueryUrl = searchURL.String()
	return searchResp, nil
}

// BookMetadata retrieves detailed metadata for a specific book from the BookBeat API
func (b *Bookbeat) BookMetadata(ctx context.Context, bookURL string) (BookResponse, error) {
	response, err := b.httpClient.R().SetContext(ctx).Get(bookURL)
	if err != nil {
		return BookResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	if !response.IsSuccessState() {
		return BookResponse{},
			fmt.Errorf("book metadata request failed with status %d: %s",
				response.StatusCode,
				response.Body)
	}

	var bookResp BookResponse
	if err := json.NewDecoder(response.Body).Decode(&bookResp); err != nil {
		return BookResponse{}, fmt.Errorf("failed to decode book metadata response: %w", err)
	}
	return bookResp, nil
}

// Search performs a comprehensive search and returns structured book data
// It searches for books, retrieves detailed metadata for each result, and transforms the data into Book structs
func (b *Bookbeat) Search(ctx context.Context, query string, author *string) ([]Book, error) {
	searchResp, err := b.SearchBooks(ctx, query, author)
	if err != nil {
		return nil, err
	}

	// Pre-allocate slice with capacity for multiple editions per book
	books := make([]Book, 0, len(searchResp.Embedded.Books)*2)
	// Process each book from search results
	for _, searchBook := range searchResp.Embedded.Books {
		// Get detailed metadata for the book
		bookResp, err := b.BookMetadata(ctx, searchBook.Links.Self.Href)
		if err != nil {
			// Skip books that fail to load metadata
			continue
		}

		// Extract series information if available
		var series *BookSeries
		if bookResp.Series != nil {
			series = &BookSeries{
				Series:   bookResp.Series.Name,
				Sequence: bookResp.Series.Partnumber,
			}
		}

		// Extract genre names
		genres := ExtractGenreNames(bookResp.Genres)

		// Clean up cover URL by removing query parameters
		cover := SanitizeCoverURL(bookResp.Cover)

		// Process each edition of the book (audiobook, ebook, etc.)
		for _, edition := range bookResp.Editions {
			// Skip unwanted formats
			if !slices.Contains(b.formats, strings.ToLower(edition.Format)) {
				continue
			}
			// Separate authors and narrators from contributors
			authors, narrators := ExtractContributors(edition.Contributors)

			// HACK: looks like ABS wants minutes and not seconds like its statet in the api definition
			if bookResp.Audiobooklength > 0 {
				bookResp.Audiobooklength = bookResp.Audiobooklength / 60
			}

			// Create book structure with all metadata
			book := Book{
				ID:            bookResp.ID,
				Title:         bookResp.Title,
				Subtitle:      bookResp.Subtitle,
				Description:   breakTagRegex.ReplaceAllString(bookResp.Summary, "\n"),
				Cover:         cover,
				Series:        series,
				Language:      bookResp.Language,
				Duration:      bookResp.Audiobooklength,
				Genres:        genres,
				Tags:          bookResp.Contenttypetags,
				Type:          edition.Format,
				Authors:       strings.Join(authors, ", "),
				Narrators:     strings.Join(narrators, ", "),
				Publisher:     edition.Publisher,
				PublishedYear: edition.Published.Year(),
				ISBN:          edition.Isbn,
			}
			books = append(books, book)
		}
	}
	return books, nil
}
