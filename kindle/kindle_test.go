package kindle_test

import (
	"context"
	"testing"

	"github.com/ahobsonsayers/abs-tract/kindle"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

const (
	TheHobbitTitle  = "The Hobbit"
	TheHobbitAuthor = "J. R. R. Tolkien"
)

func TestSearchBook(t *testing.T) {
	t.Skip("Fails in CI")
	// Should return https://www.amazon.com/dp/B007978NU6
	books, err := kindle.DefaultClient.Search(
		context.Background(),
		TheHobbitTitle,
		lo.ToPtr(TheHobbitAuthor),
	)
	require.NoError(t, err)
	require.NotEmpty(t, books)

	book := books[0]
	require.Equal(t, "B007978NU6", book.ASIN)
	require.Equal(t, "The Hobbit: 75th Anniversary Edition", book.Title)
	require.Equal(t, "J.R.R. Tolkien and Christopher Tolkien", book.Author)
	require.Equal(t, "https://m.media-amazon.com/images/I/61Ng-W9EhBL.jpg", book.Cover)
}
