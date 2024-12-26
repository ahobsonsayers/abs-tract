package goodreads_test

import (
	"context"
	"testing"

	"github.com/ahobsonsayers/abs-tract/goodreads"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

const (
	GameOfThronesBookId     = "13496"
	GameOfThronesBookTitle  = "A Game of Thrones"
	GameOfThronesBookAuthor = "George R.R. Martin"
)

func TestGetBookById(t *testing.T) {
	book, err := goodreads.DefaultClient.GetBookById(context.Background(), GameOfThronesBookId)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, book)
}

func TestGetBookByTitle(t *testing.T) {
	book, err := goodreads.DefaultClient.GetBookByTitle(context.Background(), GameOfThronesBookTitle, nil)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, book)
}

func TestSearchTitle(t *testing.T) {
	books, err := goodreads.DefaultClient.SearchBooks(context.Background(), GameOfThronesBookTitle, nil)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, books[0])
}

func TestSearchTitleAndAuthor(t *testing.T) {
	books, err := goodreads.DefaultClient.SearchBooks(
		context.Background(),
		GameOfThronesBookTitle,
		lo.ToPtr(GameOfThronesBookAuthor),
	)
	require.NoError(t, err)
	checkTheHobbitBookDetails(t, books[0])
}

func checkTheHobbitBookDetails(t *testing.T, book goodreads.Book) {
	require.Equal(t, GameOfThronesBookTitle, book.Work.Title())
	require.Equal(t, GameOfThronesBookId, book.BestEdition.Id)
	require.Regexp(t, "1562726234l/13496.jpg$", book.BestEdition.ImageURL)
	require.Equal(t, "English", book.BestEdition.Language)
	require.Equal(t, GameOfThronesBookAuthor, book.Authors[0].Name)
	require.Equal(t, "A Song of Ice and Fire", book.Series[0].Series.Title)
	require.Equal(t, "1", *book.Series[0].BookPosition)
	require.EqualValues(t, []string{"Fantasy", "Fiction", "High Fantasy"}, book.Genres)
}
