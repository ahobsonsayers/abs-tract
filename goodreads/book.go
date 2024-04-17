package goodreads

import (
	"encoding/xml"
	"math"
	"strings"
)

type Book struct {
	Work        Work            `xml:"work"`
	BestEdition Edition         // Unmarshalled using the custom unmarshaler below
	Authors     []AuthorDetails `xml:"authors>author"`
	Series      []SeriesBook    `xml:"series_works>series_work"`
	Genres      Genres          `xml:"popular_shelves"` // The (max) first 5 "genre" shelves
}

func (b *Book) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type bookAlias Book
	type bookAux struct {
		bookAlias
		Edition
	}

	var book bookAux
	err := d.DecodeElement(&book, &start)
	if err != nil {
		return err
	}

	*b = Book(book.bookAlias)
	b.BestEdition = book.Edition
	return nil
}

type Work struct {
	Title         string `xml:"original_title"`
	MediaType     string `xml:"media_type"`
	EditionsCount int    `xml:"books_count"`
	Language      int    `xml:"original_language_id"`

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
	Id               string  `xml:"id"` // Can be used to show book
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
	LanguageCode     string  `xml:"language_code"`
}

type SeriesBook struct {
	Series       Series `xml:"series"`
	BookPosition *int   `xml:"user_position"`
}

type Series struct {
	Id               string `xml:"id"`
	Title            string `xml:"title"`
	Description      string `xml:"description"`
	PrimaryBookCount int    `xml:"primary_work_count"`
	TotalBookCount   int    `xml:"series_works_count"`
	Numbered         bool   `xml:"numbered"`
}

func (s *Series) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type seriesAlias Series
	var series seriesAlias
	err := d.DecodeElement(&series, &start)
	if err != nil {
		return err
	}
	*s = Series(series)

	// Cleanup some fields
	s.Title = strings.TrimSpace(s.Title)
	s.Description = strings.TrimSpace(s.Description)

	return nil
}
