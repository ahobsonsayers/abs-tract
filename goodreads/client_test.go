package goodreads_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ahobsonsayers/abs-goodreads/goodreads"
	"github.com/stretchr/testify/require"
)

const TheHobbitBookId = "5907"

func TestGetBook(t *testing.T) {
	client := goodreads.NewClient(http.DefaultClient, nil, nil)

	book, err := client.GetBook(context.Background(), TheHobbitBookId)
	require.NoError(t, err)

	require.Equal(t, "The Hobbit", book.Work.Title)
	require.Equal(t, TheHobbitBookId, book.BestEdition.Id)
	require.Equal(t, "J.R.R. Tolkien", book.Authors[0].Name)
	require.Equal(t, "The Lord of the Rings", book.Series[0].Series.Title)
	require.Equal(t, 0, *book.Series[0].BookPosition)
}

func TestSearchBook(t *testing.T) {
	client := goodreads.NewClient(http.DefaultClient, nil, nil)

	theHobbitSearchQuery := "The Hobbit"
	books, err := client.SearchBook(context.Background(), theHobbitSearchQuery, nil)
	require.NoError(t, err)

	// Check first book returned
	book := books[0]
	require.Equal(t, "The Hobbit", book.Work.Title)
	require.Equal(t, TheHobbitBookId, book.BestEdition.Id)
	require.Equal(t, "J.R.R. Tolkien", book.Authors[0].Name)
	require.Equal(t, "The Lord of the Rings", book.Series[0].Series.Title)
	require.Equal(t, 0, *book.Series[0].BookPosition)
}
