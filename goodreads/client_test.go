package goodreads_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ahobsonsayers/audiobookshelf-goodreads/goodreads"
	"github.com/stretchr/testify/require"
)

const TheHobbitBookId = "5907"

func TestSearchBook(t *testing.T) {
	client := goodreads.NewClient(http.DefaultClient, nil, nil)

	theHobbitSearchQuery := "The Hobbit"
	works, err := client.SearchBook(context.Background(), theHobbitSearchQuery, nil)
	require.NoError(t, err)

	// Check first book returned
	book := works[0].Work
	require.Equal(t, "The Hobbit", book.Title)
	require.Equal(t, "J.R.R. Tolkien", works[0].Authors[0].Name)
}

func TestGetBook(t *testing.T) {
	client := goodreads.NewClient(http.DefaultClient, nil, nil)

	book, err := client.GetBook(context.Background(), TheHobbitBookId)
	require.NoError(t, err)

	require.Equal(t, "The Hobbit", book.Work.Title)
	require.Equal(t, TheHobbitBookId, book.BestEdition.Id)
	require.Equal(t, "J.R.R. Tolkien", book.Authors[0].Name)
}
