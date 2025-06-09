package bookbeat_test

import (
	"testing"

	"github.com/ahobsonsayers/abs-tract/bookbeat"
	"github.com/stretchr/testify/assert"
)

func TestExtractGenreNames(t *testing.T) {
	tests := []struct {
		name     string
		genres   []bookbeat.Genre
		expected []string
	}{
		{
			name: "multiple genres",
			genres: []bookbeat.Genre{
				{Genreid: 1, Name: "Fiction"},
				{Genreid: 2, Name: "Mystery"},
				{Genreid: 3, Name: "Thriller"},
			},
			expected: []string{"Fiction", "Mystery", "Thriller"},
		},
		{
			name: "single genre",
			genres: []bookbeat.Genre{
				{Genreid: 1, Name: "Non-Fiction"},
			},
			expected: []string{"Non-Fiction"},
		},
		{
			name:     "empty genres",
			genres:   []bookbeat.Genre{},
			expected: []string{},
		},
		{
			name:     "nil genres",
			genres:   nil,
			expected: []string{},
		},
		{
			name: "genres with empty names",
			genres: []bookbeat.Genre{
				{Genreid: 1, Name: ""},
				{Genreid: 2, Name: "Fantasy"},
			},
			expected: []string{"", "Fantasy"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bookbeat.ExtractGenreNames(tt.genres)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractContributorsWithMixedContributors(t *testing.T) {
	contributors := []bookbeat.Contributor{
		{ID: 1, Displayname: "John Doe", Role: "bb-author"},
		{ID: 2, Displayname: "Jane Smith", Role: "bb-narrator"},
		{ID: 3, Displayname: "Bob Wilson", Role: "bb-author"},
		{ID: 4, Displayname: "Alice Johnson", Role: "bb-narrator"},
		{ID: 5, Displayname: "Other Role", Role: "bb-translator"},
	}
	authors, narrators := bookbeat.ExtractContributors(contributors)
	assert.Equal(t, []string{"John Doe", "Bob Wilson"}, authors)
	assert.Equal(t, []string{"Jane Smith", "Alice Johnson"}, narrators)
}

func TestExtractContributorsWithOnlyAuthors(t *testing.T) {
	contributors := []bookbeat.Contributor{
		{ID: 1, Displayname: "Author One", Role: "bb-author"},
		{ID: 2, Displayname: "Author Two", Role: "bb-author"},
	}
	authors, narrators := bookbeat.ExtractContributors(contributors)
	assert.Equal(t, []string{"Author One", "Author Two"}, authors)
	assert.Equal(t, []string{}, narrators)
}

func TestExtractContributorsWithOnlyNarrators(t *testing.T) {
	contributors := []bookbeat.Contributor{
		{ID: 1, Displayname: "Narrator One", Role: "bb-narrator"},
		{ID: 2, Displayname: "Narrator Two", Role: "bb-narrator"},
	}
	authors, narrators := bookbeat.ExtractContributors(contributors)
	assert.Equal(t, []string{}, authors)
	assert.Equal(t, []string{"Narrator One", "Narrator Two"}, narrators)
}

func TestExtractContributorsWithNoRelevantRoles(t *testing.T) {
	contributors := []bookbeat.Contributor{
		{ID: 1, Displayname: "Translator", Role: "bb-translator"},
		{ID: 2, Displayname: "Editor", Role: "bb-editor"},
		{ID: 3, Displayname: "Unknown", Role: "unknown-role"},
	}
	authors, narrators := bookbeat.ExtractContributors(contributors)
	assert.Equal(t, []string{}, authors)
	assert.Equal(t, []string{}, narrators)
}

func TestExtractContributorsWithEmptyContributors(t *testing.T) {
	authors, narrators := bookbeat.ExtractContributors([]bookbeat.Contributor{})
	assert.Equal(t, []string{}, authors)
	assert.Equal(t, []string{}, narrators)
}
func TestExtractContributorsWithNilContributors(t *testing.T) {
	authors, narrators := bookbeat.ExtractContributors(nil)
	assert.Equal(t, []string{}, authors)
	assert.Equal(t, []string{}, narrators)
}

func TestExtractContributorsWithEmptyNames(t *testing.T) {
	contributors := []bookbeat.Contributor{
		{ID: 1, Displayname: "", Role: "bb-author"},
		{ID: 2, Displayname: "Valid Author", Role: "bb-author"},
		{ID: 3, Displayname: "", Role: "bb-narrator"},
	}
	authors, narrators := bookbeat.ExtractContributors(contributors)
	assert.Equal(t, []string{"Valid Author"}, authors)
	assert.Equal(t, []string{}, narrators)
}

func TestSanitizeCoverURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with query parameters",
			input:    "https://example.com/cover.jpg?param=value&size=large",
			expected: "https://example.com/cover.jpg",
		},
		{
			name:     "URL without query parameters",
			input:    "https://example.com/cover.jpg",
			expected: "https://example.com/cover.jpg",
		},
		{
			name:     "URL with fragment",
			input:    "https://example.com/cover.jpg?size=large#fragment",
			expected: "https://example.com/cover.jpg",
		},
		{
			name:     "Invalid URL",
			input:    "not-a-valid-url",
			expected: "not-a-valid-url",
		},
		{
			name:     "Empty URL",
			input:    "",
			expected: "",
		},
		{
			name:     "Error URL",
			input:    ":",
			expected: ":",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bookbeat.SanitizeCoverURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
