package goodreads

import (
	"regexp"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"golang.org/x/exp/slices"
)

func sortBookByTitleSimilarity(books []Book, title string) {
	normalisedDesiredTitle := normaliseBookTitle(title)

	similarities := make(map[string]float64)
	for _, book := range books {
		similarities[book.Work.FullTitle] = strutil.Similarity(
			normaliseBookTitle(book.Work.FullTitle),
			normalisedDesiredTitle,
			metrics.NewJaroWinkler(),
		)
	}

	slices.SortStableFunc(books, func(i, j Book) bool {
		return similarities[i.Work.FullTitle] > similarities[j.Work.FullTitle]
	})
}

var (
	spaceRegex        = regexp.MustCompile(`\s+`)
	alphanumericRegex = regexp.MustCompile(`[^a-z0-9]`)
)

func normaliseBookTitle(title string) string {
	// Normalise and trim whitespace
	title = spaceRegex.ReplaceAllString(title, " ")
	title = strings.TrimSpace(title)

	// Normalise and trim text
	title = strings.ToLower(title)
	title = alphanumericRegex.ReplaceAllString(title, "")
	title = strings.TrimPrefix(title, "the ")

	return title
}
