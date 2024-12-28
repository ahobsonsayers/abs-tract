package goodreads_test

import (
	"context"
	"testing"

	"github.com/ahobsonsayers/abs-tract/goodreads"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

const (
	TheHobbitId     = "5907"
	TheHobbitTitle  = "The Hobbit, or There and Back Again"
	TheHobbitAuthor = "J.R.R. Tolkien"
)

func TestGetBookById(t *testing.T) {
	book, err := goodreads.DefaultClient.GetBookById(context.Background(), TheHobbitId)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, book)
}

func TestGetBookByTitle(t *testing.T) {
	book, err := goodreads.DefaultClient.GetBookByTitle(context.Background(), TheHobbitTitle, nil)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, book)
}

func TestSearchTitle(t *testing.T) {
	books, err := goodreads.DefaultClient.SearchBooks(context.Background(), TheHobbitTitle, nil)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, books[0])
}

func TestSearchTitleAndAuthor(t *testing.T) {
	books, err := goodreads.DefaultClient.SearchBooks(
		context.Background(),
		TheHobbitTitle,
		lo.ToPtr(TheHobbitAuthor),
	)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, books[0])
}

func checkTheHobbitBookDetails(t *testing.T, book goodreads.Book) {
	require.Equal(t, TheHobbitTitle, book.BestEdition.Title())
	require.Equal(t, TheHobbitId, book.BestEdition.Id)
	require.Regexp(t, "1546071216l/5907.jpg$", book.BestEdition.ImageURL)
	require.Equal(t, "English", book.BestEdition.Language)
	require.Equal(t, TheHobbitAuthor, book.Authors[0].Name)
	require.Equal(t, "Middle Earth", book.Series[0].Series.Title)
	require.Equal(t, "0", *book.Series[0].BookPosition)
	require.EqualValues(t, []string{"Fantasy", "Fiction", "Classic"}, book.Genres)
}
