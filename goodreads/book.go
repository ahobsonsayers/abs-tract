package goodreads

import (
	"encoding/xml"
	"math"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

type Genres []string

// None exhaustive list of popular shelves that are known to not be genres.
// Used to allow extraction of genres from goodreads shelves.
// TODO is there a better way to recognise genres?
var nonGenreShelves = mapset.NewSet(
	"audiobook",
	"audiobooks",
	"books-i-own",
	"currently-reading",
	"default",
	"favourites",
	"library",
	"literature",
	"my-books",
	"my-library",
	"novels",
	"owned-books",
	"owned",
	"re-read",
	"series",
	"to-read",
)

func (g *Genres) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Genre aux is a struct matching the xml
	var genresAux struct {
		Genres []struct {
			Name string `xml:"name,attr"`
		} `xml:"shelf"`
	}
	err := d.DecodeElement(&genresAux, &start)
	if err != nil {
		return err
	}

	genres := make(Genres, 0, len(genresAux.Genres))
	for _, genreAux := range genresAux.Genres {
		genre := genreAux.Name

		// Skip non genre shelves
		if nonGenreShelves.Contains(genre) {
			continue
		}

		// Make genre human readable
		genre = strings.ReplaceAll(genre, "-", " ")
		genreWords := strings.Fields(genre)
		for i, word := range genreWords {
			genreWords[i] = cases.Title(language.Und).String(word) // Capitalize each word
		}
		genre = strings.Join(genreWords, " ")

		genres = append(genres, genre)
		if len(genres) == 5 {
			break
		}
	}
	*g = genres

	return nil
}
