package goodreads

import (
	"encoding/xml"
	"math"
	"regexp"
	"strings"

	"github.com/ahobsonsayers/abs-tract/utils"
	"github.com/k3a/html2text"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var (
	alternativeCoverRegex = regexp.MustCompile(`^\s*<i>.*[Aa]lternat(iv)?e cover.*</i>\s*$`)
	lastBracketRegex      = regexp.MustCompile(`^(.*)(\([^\(\)]*\))([^()]*)$`)
)

type BookOverview struct {
	Id     string `xml:"id"`
	Title  string `xml:"title"`
	Author string `xml:"author>name"`
}

func (b *BookOverview) Sanitise() {
	// Strip last brackets from title. This is the series.
	// e.g. Harry Potter and the Chamber of Secrets (Harry Potter, #2)
	b.Title = lastBracketRegex.ReplaceAllString(b.Title, "$1$2")
}

type Book struct {
	Work        Work            `xml:"work"`
	BestEdition Edition         // Unmarshalled using the custom unmarshaler below
	Authors     []AuthorDetails `xml:"authors>author"`
	Series      []SeriesBook    `xml:"series_works>series_work"`
	Genres      Genres          `xml:"popular_shelves"` // The first "genre" shelves
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
	FullTitle     string `xml:"original_title"`
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

// Title is the full title with any subtitle removed.
// A subtitle is anything after the first : in the full title
func (w Work) Title() string {
	titleParts := strings.Split(w.FullTitle, ":")
	return strings.TrimSpace(titleParts[0])
}

// Subtitle is the subtle part of the full title.
// A subtitle is anything after the first : in the full title
func (w Work) Subtitle() string {
	colonIdx := strings.Index(w.FullTitle, ":")
	if colonIdx == -1 {
		return ""
	}
	return strings.TrimSpace(w.FullTitle[colonIdx+1:])
}

func (w Work) AverageRating() float64 {
	averageRating := float64(w.RatingsSum) / float64(w.RatingsCount)
	return math.Round(averageRating*100) / 100 // Round to two decimal places
}

type Edition struct {
	Id               string  `xml:"id"`
	ISBN             *string `xml:"isbn13"`
	FullTitle        string  `xml:"title"`
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

// Title is the full title with any subtitle removed.
// A subtitle is anything after the first : in the full title
func (e Edition) Title() string {
	titleParts := strings.Split(e.FullTitle, ":")
	return strings.TrimSpace(titleParts[0])
}

// Subtitle is the subtle part of the full title.
// A subtitle is anything after the first : in the full title
func (e Edition) Subtitle() string {
	colonIdx := strings.Index(e.FullTitle, ":")
	if colonIdx == -1 {
		return ""
	}
	return strings.TrimSpace(e.FullTitle[colonIdx+1:])
}

func (e *Edition) Sanitise() {
	// Description can sometimes be html and contain preamble about alternative covers
	description := strings.TrimSpace(e.Description)
	description = alternativeCoverRegex.ReplaceAllString(description, "")
	description = html2text.HTML2Text(description)
	e.Description = description

	// Get original cover image by cleaning the ul0
	if strings.Contains(e.ImageURL, "nophoto") {
		e.ImageURL = ""
	} else {
		e.ImageURL = utils.SanitiseImageURL(e.ImageURL)
	}

	// Convert language from code to name (if possible)
	lang, err := language.Parse(e.Language)
	if err == nil {
		e.Language = display.English.Languages().Name(lang)
	} else {
		e.Language = strings.ToTitle(e.Language)
	}
}

func BookIds(books []BookOverview) []string {
	ids := make([]string, 0, len(books))
	for _, book := range books {
		ids = append(ids, book.Id)
	}
	return ids
}

func BookTitles(books []BookOverview) []string {
	titles := make([]string, 0, len(books))
	for _, book := range books {
		titles = append(titles, book.Title)
	}
	return titles
}
