package goodreads

import (
	"math"
	"regexp"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"golang.org/x/exp/slices"
)

func sortBookByTitleSimilarity(books []Book, title string) {
	normalisedDesiredTitle := normaliseBookTitle(title)

	// Get the best similarity from the work and best edition titles.
	// This is useful if the tiles differ in different regions.
	similarities := make(map[string]float64)
	for _, book := range books {
		workTitleSimilarity := strutil.Similarity(
			normaliseBookTitle(book.Work.FullTitle),
			normalisedDesiredTitle,
			metrics.NewJaroWinkler(),
		)
		bestEditionTitleSimilarity := strutil.Similarity(
			normaliseBookTitle(book.BestEdition.FullTitle),
			normalisedDesiredTitle,
			metrics.NewJaroWinkler(),
		)
		similarities[book.BestEdition.Id] = math.Max(workTitleSimilarity, bestEditionTitleSimilarity)
	}

	slices.SortStableFunc(books, func(i, j Book) bool {
		return similarities[i.BestEdition.Id] > similarities[j.BestEdition.Id]
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
