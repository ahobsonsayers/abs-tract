package goodreads

import (
	"regexp"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"golang.org/x/exp/slices"
)

func sortBookOverviewsByTitleSimilarity(bookOverviews []BookOverview, title string) {
	normalisedDesiredTitle := normaliseBookTitle(title)

	similarities := make(map[string]float64)
	for _, book := range bookOverviews {
		similarities[book.Title] = strutil.Similarity(
			normaliseBookTitle(book.Title),
			normalisedDesiredTitle,
			metrics.NewJaroWinkler(),
		)
	}

	slices.SortStableFunc(bookOverviews, func(i, j BookOverview) bool {
		return similarities[i.Title] > similarities[j.Title]
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
