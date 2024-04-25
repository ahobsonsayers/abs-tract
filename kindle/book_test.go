package kindle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseBookCoversAttrValue(t *testing.T) {
	coverUrl := "https://m.media-amazon.com/images/I/61Ng-W9EhBL.jpg"

	coverSetAttrValue := "https://m.media-amazon.com/images/I/61Ng-W9EhBL._AC_UY218_.jpg 1x, https://m.media-amazon.com/images/I/61Ng-W9EhBL._AC_UY327_QL65_.jpg 1.5x, https://m.media-amazon.com/images/I/61Ng-W9EhBL._AC_UY436_QL65_.jpg 2x, https://m.media-amazon.com/images/I/61Ng-W9EhBL._AC_UY500_QL65_.jpg 2.2935x" // nolint
	parsedCoverUrl := parseBookCoversAttrValue(coverSetAttrValue)
	require.Equal(t, coverUrl, parsedCoverUrl)
}

func TestParseBookInfoNodeValue(t *testing.T) {
	author := "J.R.R. Tolkien"
	publishDate := time.Date(2012, time.February, 15, 0, 0, 0, 0, time.UTC)

	bookInfoNodeValue := "by J.R.R. Tolkien | Sold by: HarperCollins Publishers  | Feb 15, 2012"
	parsedAuthor, parsedPublishDate := parseBookInfoNodeValue(bookInfoNodeValue)
	require.Equal(t, author, parsedAuthor)
	require.Equal(t, publishDate, *parsedPublishDate)
}
