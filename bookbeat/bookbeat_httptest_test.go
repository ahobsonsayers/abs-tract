package bookbeat_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ahobsonsayers/abs-tract/bookbeat"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadTestData loads JSON test data from testdata directory
func loadTestData(t *testing.T, filename string) []byte {
	t.Helper()
	path := filepath.Join("testdata", filename)
	data, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read test data file: %s", filename)
	return data
}

func TestSearchBooksWithMockServer(t *testing.T) {
	// Load test data from JSON file
	searchResponseData := loadTestData(t, "search_response.json")

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request parameters
		assert.Equal(t, "/api/next/search", r.URL.Path)
		assert.Contains(t, r.URL.Query().Get("query"), "Harry Potter")
		assert.Equal(t, "audiobook", r.URL.Query().Get("format"))
		assert.Equal(t, "english", r.URL.Query().Get("language"))

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(searchResponseData); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with test server URLs
	client, err := bookbeat.NewClientWithURLs("ee", "audiobook", "en", server.URL+"/api/next/search")
	require.NoError(t, err)

	// Test the search
	result, err := client.SearchBooks(context.Background(), TestBookTitle, lo.ToPtr(TestBookAuthor))
	require.NoError(t, err)

	// Verify results match test data
	assert.Equal(t, 1, result.Count)
	assert.Len(t, result.Embedded.Books, 1)
	assert.Equal(t, 38059, result.Embedded.Books[0].ID)
	assert.Contains(t, result.QueryUrl, "query=Harry+Potter")
}

func TestBookMetadataWithMockServer(t *testing.T) {
	// Load test data from JSON file
	bookMetadataData := loadTestData(t, "book_metadata.json")

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/books/372/38059", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(bookMetadataData); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with test server URLs
	client, err := bookbeat.NewClientWithURLs("ee", "all", "en", server.URL+"/api/next/search")
	require.NoError(t, err)

	// Test book metadata retrieval
	result, err := client.BookMetadata(context.Background(), server.URL+"/api/books/372/38059")
	require.NoError(t, err)

	// Verify results match test data
	assert.Equal(t, 38059, result.ID)
	assert.Equal(t, TestBookTitle, result.Title)
	assert.Empty(t, result.Subtitle) // null in JSON
	assert.Contains(t, result.Summary, "<br>")
	assert.Equal(t, "English", result.Language)
	assert.Equal(t, 30334, result.Audiobooklength)
	assert.Len(t, result.Genres, 8)
	assert.Equal(t, "Fantasy", result.Genres[0].Name)
	assert.NotNil(t, result.Series)
	assert.Equal(t, "Harry Potter", result.Series.Name)
	assert.Equal(t, 1, result.Series.Partnumber)
	assert.Len(t, result.Editions, 2) // Both audiobook and ebook
	assert.Equal(t, "audioBook", result.Editions[0].Format)
	assert.Equal(t, "eBook", result.Editions[1].Format)
}

func TestSearchWithMockServer(t *testing.T) {
	// Load test data from JSON files
	searchResponseData := loadTestData(t, "search_response.json")
	bookMetadataData := loadTestData(t, "book_metadata.json")

	// Create test server that handles both endpoints
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch path := r.URL.Path; path {
		case "/api/next/search":
			// Modify the search response to point to our test server
			var searchResp bookbeat.SearchResponse

			if err := json.Unmarshal(searchResponseData, &searchResp); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
			searchResp.Embedded.Books[0].Links.Self.Href = server.URL + "/api/books/372/38059"
			if err := json.NewEncoder(w).Encode(searchResp); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		case "/api/books/372/38059":
			if _, err := w.Write(bookMetadataData); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client with test server URLs
	client, err := bookbeat.NewClientWithURLs("ee", "audiobook", "en", server.URL+"/api/next/search")
	require.NoError(t, err)

	// Test full search functionality
	books, err := client.Search(context.Background(), TestBookTitle, nil)
	require.NoError(t, err)

	// Verify results using real test data
	require.Len(t, books, 1) // Only audiobook edition since we filtered for audiobook
	book := books[0]

	assert.Equal(t, 38059, book.ID)
	assert.Equal(t, TestBookTitle, book.Title)
	assert.Empty(t, book.Subtitle) // null in test data
	assert.Contains(t, book.Description, "Stephen Fry brings")
	assert.NotContains(t, book.Description, "<br>") // Should be converted to \n
	assert.NotContains(t, book.Cover, "?")          // Query params should be removed
	assert.Equal(t, "English", book.Language)
	assert.Equal(t, 505, book.Duration)
	assert.Contains(t, book.Genres, "Fantasy")
	assert.Contains(t, book.Tags, "contentsynch")
	assert.Equal(t, "audioBook", book.Type)
	assert.Equal(t, TestBookAuthor, book.Authors)
	assert.Equal(t, "Stephen Fry", book.Narrators)
	assert.Equal(t, "Pottermore Publishing", book.Publisher)
	assert.Equal(t, 2015, book.PublishedYear)
	assert.Equal(t, "9781781102367", book.ISBN)
}

func TestSearchBooksEmptyResults(t *testing.T) {
	// Load empty search response test data
	emptySearchData := loadTestData(t, "search_response_empty.json")

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(emptySearchData); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with test server URLs
	client, err := bookbeat.NewClientWithURLs("ee", "all", "en", server.URL+"/api/next/search")
	require.NoError(t, err)

	// Test search with no results
	result, err := client.SearchBooks(context.Background(), "NonexistentBook", nil)
	require.NoError(t, err)

	// Verify empty results
	assert.Equal(t, 0, result.Count)
	assert.Empty(t, result.Embedded.Books)
	assert.False(t, result.IsCapped)
}

func TestSearchBooksErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedError string
	}{
		{
			name:          "HTTP 404 error",
			statusCode:    404,
			responseBody:  "Not Found",
			expectedError: "search request failed with status 404",
		},
		{
			name:          "HTTP 500 error",
			statusCode:    500,
			responseBody:  "Internal Server Error",
			expectedError: "search request failed with status 500",
		},
		{
			name:          "Invalid JSON response",
			statusCode:    200,
			responseBody:  "invalid json",
			expectedError: "failed to decode search response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				if _, err := w.Write([]byte(tt.responseBody)); err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			}))
			defer server.Close()

			client, err := bookbeat.NewClientWithURLs("uk", "all", "en", server.URL)
			require.NoError(t, err)

			_, err = client.SearchBooks(context.Background(), "test", nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestBookMetadataErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedError string
	}{
		{
			name:          "HTTP 404 error",
			statusCode:    404,
			responseBody:  "Not Found",
			expectedError: "book metadata request failed with status 404",
		},
		{
			name:          "Invalid JSON response",
			statusCode:    200,
			responseBody:  "invalid json",
			expectedError: "failed to decode book metadata response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				if _, err := w.Write([]byte(tt.responseBody)); err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			}))
			defer server.Close()

			client, err := bookbeat.NewClientWithURLs("uk", "all", "en", server.URL)
			require.NoError(t, err)

			_, err = client.BookMetadata(context.Background(), server.URL+"/book/123")
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestSearchBooksErrors(t *testing.T) {
	client, err := bookbeat.NewClientWithURLs("uk", "all", "en", ";;")
	require.NoError(t, err)
	_, err = client.SearchBooks(context.Background(), "test", nil)
	require.Error(t, err)
}

func TestBookMetadataErrors(t *testing.T) {
	client, err := bookbeat.NewClientWithURLs("uk", "all", "en", ";;")
	require.NoError(t, err)
	_, err = client.BookMetadata(context.Background(), ";;")
	require.Error(t, err)
}

func TestSearchErrors(t *testing.T) {
	client, err := bookbeat.NewClientWithURLs("uk", "all", "en", ";;")
	require.NoError(t, err)
	_, err = client.Search(context.Background(), "test", nil)
	require.Error(t, err)
}

func TestSearchWithMetadataError(t *testing.T) {
	// Load test data from JSON files
	searchResponseData := loadTestData(t, "search_response.json")
	bookMetadataData := loadTestData(t, "book_metadata.json")

	// Create test server that handles both endpoints
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch path := r.URL.Path; path {
		case "/api/next/search":
			// Modify the search response to point to our test server
			var searchResp bookbeat.SearchResponse

			if err := json.Unmarshal(searchResponseData, &searchResp); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
			searchResp.Embedded.Books[0].Links.Self.Href = ";;" // Set invalid URL
			if err := json.NewEncoder(w).Encode(searchResp); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		case "/api/books/372/38059":
			if _, err := w.Write(bookMetadataData); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client with test server URLs
	client, err := bookbeat.NewClientWithURLs("ee", "audiobook", "en", server.URL+"/api/next/search")
	require.NoError(t, err)

	books, err := client.Search(context.Background(), TestBookTitle, nil)
	require.NoError(t, err)
	require.Empty(t, books)
}
