package bookbeat

import (
	"net/url"
)

// ExtractGenreNames selects only the name from slice of genres
func ExtractGenreNames(genres []Genre) []string {
	list := make([]string, 0, len(genres))
	for _, genre := range genres {
		list = append(list, genre.Name)
	}
	return list
}

// SanitizeCoverURL removes query parameters and fragments from cover URLs
func SanitizeCoverURL(coverURL string) string {
	if u, err := url.Parse(coverURL); err == nil {
		u.RawQuery = ""
		u.Fragment = ""
		return u.String()
	}
	return coverURL
}

// ExtractContributors separates authors and narrators from contributors
func ExtractContributors(contributors []Contributor) (authors, narrators []string) {
	authors = make([]string, 0, len(contributors))
	narrators = make([]string, 0, len(contributors))

	for _, entry := range contributors {
		if entry.Displayname == "" {
			continue
		}
		switch entry.Role {
		case "bb-author":
			authors = append(authors, entry.Displayname)
		case "bb-narrator":
			narrators = append(narrators, entry.Displayname)
		}
	}
	return authors, narrators
}
