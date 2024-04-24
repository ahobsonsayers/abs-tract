package kindle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBookAuthorFromInfoNodeValue(t *testing.T) {
	authorNodeValue := "by J.R.R. Tolkien | Sold by: HarperCollins Publishers  | Feb 15, 2012"
	author := bookAuthorFromInfoNodeValue(authorNodeValue)
	expectedAuthor := "J.R.R. Tolkien"
	require.Equal(t, expectedAuthor, author)
}

func TestBookCoverFromCoverSetNodeValue(t *testing.T) {
	coverSetNodeValue := "https://m.media-amazon.com/images/I/61Ng-W9EhBL._AC_UY218_.jpg 1x, https://m.media-amazon.com/images/I/61Ng-W9EhBL._AC_UY327_QL65_.jpg 1.5x, https://m.media-amazon.com/images/I/61Ng-W9EhBL._AC_UY436_QL65_.jpg 2x, https://m.media-amazon.com/images/I/61Ng-W9EhBL._AC_UY500_QL65_.jpg 2.2935x" // nolint
	cover := bookCoverFromCoverSetAttrValue(coverSetNodeValue)
	expectedCover := "https://m.media-amazon.com/images/I/61Ng-W9EhBL._AC_UY500_QL65_.jpg"
	require.Equal(t, expectedCover, cover)
}
