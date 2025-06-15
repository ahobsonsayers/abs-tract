package bookbeat_test

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/ahobsonsayers/abs-tract/bookbeat"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

const (
	TestBookTitle  = "Harry Potter and the Philosopher's Stone"
	TestBookAuthor = "J.K. Rowling"
	BookURL        = "https://edge.bookbeat.com/api/books/372/38059"
)

func integrationTest(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests, set environment variable INTEGRATION")
	}
}

func TestSearchBooks(t *testing.T) {
	integrationTest(t)

	client, err := bookbeat.NewClient("ee", "all", "en")
	require.NoError(t, err)
	require.NotNil(t, client)

	searchResp, err := client.SearchBooks(context.Background(), TestBookTitle, nil)
	require.NoError(t, err)

	// Verify response structure
	require.NotEmpty(t, searchResp.Embedded.Books)
	require.Positive(t, searchResp.Count)

	// Check first book has required fields
	firstBook := searchResp.Embedded.Books[0]
	require.NotZero(t, firstBook.ID)
	require.NotEmpty(t, firstBook.Links.Self.Href)

	// Verify query URL is set
	encodedTitle := url.QueryEscape(TestBookTitle)
	require.NotEmpty(t, searchResp.QueryUrl)
	require.Contains(t, searchResp.QueryUrl, encodedTitle)
}

func TestSearchBooksWithAuthor(t *testing.T) {
	integrationTest(t)

	client, err := bookbeat.NewClient("ee", "all", "en")
	require.NoError(t, err)

	searchResp, err := client.SearchBooks(context.Background(), TestBookTitle, lo.ToPtr(TestBookAuthor))
	require.NoError(t, err)

	require.NotEmpty(t, searchResp.Embedded.Books)
	require.Positive(t, searchResp.Count)

	// Verify query URL contains both title and author
	encodedTitle := url.QueryEscape(TestBookTitle)
	encodedAuthor := url.QueryEscape(TestBookAuthor)
	require.Contains(t, searchResp.QueryUrl, encodedTitle)
	require.Contains(t, searchResp.QueryUrl, encodedAuthor)
}

func TestSearchBooksWithFormats(t *testing.T) {
	integrationTest(t)

	tests := []struct {
		name          string
		market        string
		format        string
		shouldBeEmpty bool
	}{
		{"audiobook only (uk)", "uk", "audiobook", false},
		{"ebook only (uk)", "uk", "ebook", true},
		{"all formats (uk)", "uk", "all", false},
		{"audiobook only (ee)", "ee", "audiobook", false},
		{"ebook only (ee)", "ee", "ebook", false},
		{"all formats (ee)", "ee", "all", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := bookbeat.NewClient(tt.market, tt.format, "en")
			require.NoError(t, err)

			searchResp, err := client.SearchBooks(context.Background(), TestBookTitle, nil)
			require.NoError(t, err)
			if tt.shouldBeEmpty {
				require.Empty(t, searchResp.Embedded.Books)
			} else {
				require.NotEmpty(t, searchResp.Embedded.Books)
			}

			// Verify format parameter is in query URL
			if tt.format != "all" {
				require.Contains(t, searchResp.QueryUrl, "format="+tt.format)
				require.Equal(t, 1, strings.Count(searchResp.QueryUrl, "format"))
			} else {
				require.Contains(t, searchResp.QueryUrl, "format=audiobook")
				require.Contains(t, searchResp.QueryUrl, "format=ebook")
				require.Equal(t, 2, strings.Count(searchResp.QueryUrl, "format"))
			}
		})
	}
}

func TestSearchWithFormats(t *testing.T) {
	integrationTest(t)

	tests := []struct {
		name     string
		market   string
		format   string
		shouldBe int
	}{
		{"audiobook only (uk)", "uk", "audiobook", 1},
		{"ebook only (uk)", "uk", "ebook", 0},
		{"all formats (uk)", "uk", "all", 1},
		{"audiobook only (ee)", "ee", "audiobook", 1},
		{"ebook only (ee)", "ee", "ebook", 1},
		{"all formats (ee)", "ee", "all", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := bookbeat.NewClient(tt.market, tt.format, "en")
			require.NoError(t, err)

			searchResp, err := client.Search(context.Background(), TestBookTitle, nil)
			require.NoError(t, err)
			require.Len(t, searchResp, tt.shouldBe)
		})
	}
}

func TestSearchBooksWithLanguages(t *testing.T) {
	integrationTest(t)

	tests := []struct {
		name      string
		languages string
	}{
		{"english only", "en"},
		{"multiple languages", "en,de"},
		{"all languages", "all"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := bookbeat.NewClient("ee", "all", tt.languages)
			require.NoError(t, err)

			searchResp, err := client.SearchBooks(context.Background(), TestBookTitle, nil)
			require.NoError(t, err)
			require.NotEmpty(t, searchResp.Embedded.Books)

			// Verify language parameter is in query URL
			switch languageCode := tt.languages; languageCode {
			case "en":
				require.Contains(t, searchResp.QueryUrl, "language=english")
				require.Equal(t, 1, strings.Count(searchResp.QueryUrl, "language"))
			case "en,de":
				require.Contains(t, searchResp.QueryUrl, "language=english")
				require.Contains(t, searchResp.QueryUrl, "language=german")
				require.Equal(t, 2, strings.Count(searchResp.QueryUrl, "language"))
			case "all":
				require.Contains(t, searchResp.QueryUrl, "language=english")
				require.Contains(t, searchResp.QueryUrl, "language=german")
				require.Greater(t, strings.Count(searchResp.QueryUrl, "language"), 10)
			}
		})
	}
}

func TestBookMetadata(t *testing.T) {
	integrationTest(t)

	// First get a book URL from search
	client, err := bookbeat.NewClient("ee", "all", "en")
	require.NoError(t, err)

	searchResp, err := client.SearchBooks(context.Background(), TestBookTitle, nil)
	require.NoError(t, err)
	require.NotEmpty(t, searchResp.Embedded.Books)

	bookURL := searchResp.Embedded.Books[0].Links.Self.Href
	require.Equal(t, "https://edge.bookbeat.com/api/books/372/38059", bookURL)

	// Get detailed metadata
	bookResp, err := client.BookMetadata(context.Background(), bookURL)
	require.NoError(t, err)

	// Verify response fields
	require.Equal(t, 38059, bookResp.ID)
	require.Equal(t, TestBookTitle, bookResp.Title)
	require.Contains(t, bookResp.Cover, "book-covers")
	require.Equal(t, "English", bookResp.Language)
	require.NotEmpty(t, bookResp.Editions)

	// Check first edition fields
	firstEdition := bookResp.Editions[0]
	require.Equal(t, 56072, firstEdition.ID)
	require.Equal(t, "audioBook", firstEdition.Format)
	require.Equal(t, "Pottermore Publishing", firstEdition.Publisher)
	require.Len(t, firstEdition.Contributors, 2)
	authors, narrators := bookbeat.ExtractContributors(firstEdition.Contributors)
	require.Contains(t, authors, "J.K. Rowling")
	require.Contains(t, narrators, "Stephen Fry")
}

func TestSearch(t *testing.T) {
	integrationTest(t)

	client, err := bookbeat.NewClient("ee", "all", "en")
	require.NoError(t, err)

	books, err := client.Search(context.Background(), TestBookTitle, nil)
	require.NoError(t, err)
	require.NotEmpty(t, books)

	// Check first book has all expected fields
	firstBook := books[0]
	require.NotZero(t, firstBook.ID)
	require.NotEmpty(t, firstBook.Title)
	require.NotEmpty(t, firstBook.Type)
	require.NotEmpty(t, firstBook.Authors)
	require.NotEmpty(t, firstBook.Cover)
	require.NotEmpty(t, firstBook.Language)
	require.NotZero(t, firstBook.PublishedYear)

	// Verify description has HTML break tags converted to newlines
	require.NotContains(t, firstBook.Description, "<br")

	// Verify cover URL has been sanitized (no query parameters)
	require.NotContains(t, firstBook.Cover, "?")
}

func TestSearchWithAuthor(t *testing.T) {
	integrationTest(t)

	client, err := bookbeat.NewClient("ee", "all", "en")
	require.NoError(t, err)

	books, err := client.Search(context.Background(), TestBookTitle, lo.ToPtr(TestBookAuthor))
	require.NoError(t, err)
	require.NotEmpty(t, books)

	firstBook := books[0]
	require.NotEmpty(t, firstBook.Authors)
	require.Contains(t, firstBook.Authors, "Rowling")
}
