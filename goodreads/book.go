package goodreads

import (
	"encoding/xml"
	"math"
	"regexp"
	"strings"

	"github.com/k3a/html2text"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var (
	alternativeCoverRegex = regexp.MustCompile(`^\s*<i>.*?[Aa]lternate cover.*?<\/i>\s*`)
	imageUrlRegex         = regexp.MustCompile(`(\d+)\..*?\.(jpe?g)`)
)

type Book struct {
	Work        Work            `xml:"work"`
	BestEdition Edition         // Unmarshalled using the custom unmarshaler below
	Authors     []AuthorDetails `xml:"authors>author"`
	Series      []SeriesBook    `xml:"series_works>series_work"`
	Genres      Genres          `xml:"popular_shelves"` // The (max) first 5 "genre" shelves
}

func (b *Book) Sanitise() {
	b.BestEdition.Sanitise()
	for idx, series := range b.Series {
		series.Sanitise()
		b.Series[idx] = series
	}
}

func (b *Book) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type alias Book
	var unmarshaller struct {
		alias
		Edition
	}
	err := d.DecodeElement(&unmarshaller, &start)
	if err != nil {
		return err
	}

	*b = Book(unmarshaller.alias)
	b.BestEdition = unmarshaller.Edition

	b.Sanitise()

	return nil
}

type Work struct {
	Title         string `xml:"original_title"`
	MediaType     string `xml:"media_type"`
	EditionsCount int    `xml:"books_count"`

	// Publication
	PublicationYear  int `xml:"original_publication_year"`
	PublicationMonth int `xml:"original_publication_month"`
	PublicationDay   int `xml:"original_publication_day"`

	// Ratings
	RatingsSum         int    `xml:"ratings_sum"`
	RatingsCount       int    `xml:"ratings_count"`
	ReviewsCount       int    `xml:"text_reviews_count"`
	RatingDistribution string `xml:"rating_dist"`
}

func (w Work) AverageRating() float64 {
	averageRating := float64(w.RatingsSum) / float64(w.RatingsCount)
	return math.Round(averageRating*100) / 100 // Round to two decimal places
}

type Edition struct {
	Id               string  `xml:"id"`
	ISBN             *string `xml:"isbn13"`
	Title            string  `xml:"title"`
	Description      string  `xml:"description"`
	NumPages         string  `xml:"num_pages"`
	ImageURL         string  `xml:"image_url"`
	URL              string  `xml:"url"`
	Format           string  `xml:"format"`
	PublicationYear  string  `xml:"publication_year"`
	PublicationMonth string  `xml:"publication_month"`
	PublicationDay   string  `xml:"publication_day"`
	Publisher        string  `xml:"publisher"`
	CountryCode      string  `xml:"country_code"`
	Language         string  `xml:"language_code"`
}

func (e *Edition) Sanitise() {
	// Description can sometimes be html and contain preamble about alternative covers
	description := strings.TrimSpace(e.Description)
	description = alternativeCoverRegex.ReplaceAllString(description, "")
	description = html2text.HTML2Text(description)
	e.Description = description

	// Get largest image by removing anything between the last number and the extensions
	// For Example:
	// https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1546071216l/5907._SX98_.jpg"
	// Should be:
	// "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1546071216l/5907.jpg"
	e.ImageURL = imageUrlRegex.ReplaceAllString(e.ImageURL, "$1.$2")

	// Convert language from code to name (if possible)
	lang, err := language.Parse(e.Language)
	if err == nil {
		e.Language = display.English.Languages().Name(lang)
	} else {
		e.Language = strings.ToTitle(e.Language)
	}
}

type SeriesBook struct {
	Series       Series  `xml:"series"`
	BookPosition *string `xml:"user_position"`
}

func (s *SeriesBook) Sanitise() {
	s.Series.Sanitise()
}

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
