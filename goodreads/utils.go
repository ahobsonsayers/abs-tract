package goodreads

import (
	"regexp"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"golang.org/x/exp/slices"
)

func sortBookByAuthorSimilarity(books []Book, author string) {
	normalisedDesiredAuthor := normaliseString(author)

	// Get the best similarity of all authors of the book
	authorSimilarities := make(map[string]float64)
	for _, book := range books {
		bestAuthorSimilarity := 0.0
		for _, author := range book.Authors {
			authorSimilarity := strutil.Similarity(
				normaliseString(author.Name),
				normalisedDesiredAuthor,
				metrics.NewJaroWinkler(),
			)
			if authorSimilarity > bestAuthorSimilarity {
				bestAuthorSimilarity = authorSimilarity
			}
		}
		authorSimilarities[book.BestEdition.Id] = bestAuthorSimilarity
	}

	slices.SortStableFunc(books, func(i, j Book) bool {
		return authorSimilarities[i.BestEdition.Id] > authorSimilarities[j.BestEdition.Id]
	})
}

var (
	spaceRegex        = regexp.MustCompile(`\s+`)
	alphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]`)
)

// normaliseString normalises a string (e.g. a title or author) by:
// - Replacing (normalising) all whitespace with a single space character
// - Removing any leading or training whitespace
// - Removing any non alpha, numerical or space characters
// - Converting text to lowercase
// - Removing the "the " prefix
func normaliseString(s string) string {
	// Normalise and trim whitespace
	s = spaceRegex.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)

	// Normalise and trim text
	s = alphanumericRegex.ReplaceAllString(s, "")
	s = strings.ToLower(s)
	s = strings.TrimPrefix(s, "the ")

	return s
}
