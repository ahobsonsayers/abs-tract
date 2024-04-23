package goodreads

import (
	"strings"
)

type Series struct {
	Id               string `xml:"id"`
	Title            string `xml:"title"`
	Description      string `xml:"description"`
	PrimaryBookCount int    `xml:"primary_work_count"`
	TotalBookCount   int    `xml:"series_works_count"`
	Numbered         bool   `xml:"numbered"`
}

func (s *Series) Sanitise() {
	s.Title = strings.TrimSpace(s.Title)
	s.Description = strings.TrimSpace(s.Description)
}

type SeriesBook struct {
	Series       Series  `xml:"series"`
	BookPosition *string `xml:"user_position"`
}

func (s *SeriesBook) Sanitise() {
	s.Series.Sanitise()
}
