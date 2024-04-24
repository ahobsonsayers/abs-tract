package goodreads_test

import (
	"context"
	"testing"

	"github.com/ahobsonsayers/abs-goodreads/goodreads"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

const (
	TheHobbitBookId     = "5907"
	TheHobbitBookTitle  = "The Hobbit"
	TheHobbitBookAuthor = "J.R.R. Tolkien"
)

func TestGetBookById(t *testing.T) {
	book, err := goodreads.DefaultClient.GetBookById(context.Background(), TheHobbitBookId)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, book)
}

func TestGetBookByTitle(t *testing.T) {
	book, err := goodreads.DefaultClient.GetBookByTitle(context.Background(), TheHobbitBookTitle, nil)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, book)
}

func TestSearchTitle(t *testing.T) {
	books, err := goodreads.DefaultClient.SearchBooks(context.Background(), TheHobbitBookTitle, nil)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, books[0])
}

func TestSearchTitleAndAuthor(t *testing.T) {
	books, err := goodreads.DefaultClient.SearchBooks(
		context.Background(),
		TheHobbitBookTitle,
		lo.ToPtr(TheHobbitBookAuthor),
	)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, books[0])
}

func checkTheHobbitBookDetails(t *testing.T, book goodreads.Book) {
	require.Equal(t, TheHobbitBookTitle, book.Work.Title)
	require.Equal(t, TheHobbitBookId, book.BestEdition.Id)
	require.Regexp(t, "1546071216l/5907.jpg$", book.BestEdition.ImageURL)
	require.Equal(t, "English", book.BestEdition.Language)
	require.Equal(t, TheHobbitBookAuthor, book.Authors[0].Name)
	require.Equal(t, "The Lord of the Rings", book.Series[0].Series.Title)
	require.Equal(t, "0", *book.Series[0].BookPosition)
	require.EqualValues(t, []string{"Fantasy", "Classic", "Fiction"}, book.Genres)
}
