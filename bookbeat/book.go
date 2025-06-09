package bookbeat

type BookSeries struct {
	Series   string `json:"series"`
	Sequence int    `json:"sequence,omitempty"`
}

type Book struct {
	ID            int         `json:"id"`
	Type          string      `json:"type"`
	Title         string      `json:"title"`
	Subtitle      string      `json:"subtitle,omitempty"`
	Authors       string      `json:"authors"`
	Narrators     string      `json:"narrators,omitempty"`
	Publisher     string      `json:"publisher,omitempty"`
	PublishedYear int         `json:"publishedYear,omitempty"`
	Description   string      `json:"description,omitempty"`
	ISBN          string      `json:"isbn,omitempty"`
	ASIN          string      `json:"asin,omitempty"`
	Language      string      `json:"language,omitempty"`
	Duration      int         `json:"duration,omitempty"`
	Genres        []string    `json:"genres,omitempty"`
	Tags          []string    `json:"tags,omitempty"`
	Series        *BookSeries `json:"series,omitempty"`
	Cover         string      `json:"cover,omitempty"`
}
