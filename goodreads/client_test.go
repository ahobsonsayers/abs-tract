package goodreads_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ahobsonsayers/abs-goodreads/goodreads"
	"github.com/ahobsonsayers/abs-goodreads/utils"
	"github.com/stretchr/testify/require"
)

const (
	TheHobbitBookId     = "5907"
	TheHobbitBookTitle  = "The Hobbit"
	TheHobbitBookAuthor = "J.R.R. Tolkien"
)

func TestGetBookById(t *testing.T) {
	client := goodreads.NewClient(http.DefaultClient, nil, nil)

	book, err := client.GetBookById(context.Background(), TheHobbitBookId)
	require.NoError(t, err)

	checkTheHobbitBookDetails(t, book)
}

func TestGetBookByTitle(t *testing.T) {
	client := goodreads.NewClient(http.DefaultClient, nil, nil)

	book, err := client.GetBookByTitle(context.Background(), TheHobbitBookTitle, nil)
	require.NoError(t, err)

	checkTheHobbitBookDetails(t, book)
}

func TestSearch(t *testing.T) {
	client := goodreads.NewClient(http.DefaultClient, nil, nil)

	books, err := client.SearchBooks(
		context.Background(),
		TheHobbitBookTitle,
		utils.ToPointer(TheHobbitBookAuthor),
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
