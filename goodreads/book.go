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
	// These are dirty workarounds, but they seem to work

	// Regex to match last brackets from title. This is the series.
	// e.g. Harry Potter and the Chamber of Secrets (Harry Potter, #2)
	titleSeriesRegex = regexp.MustCompile(`\([^)]*#\d+(\.\d+)?\)$`)

	// Regex to match alternative cover preamble in description.
	// e.g. Harry Potter and the Chamber of Secrets (Harry Potter, #2)
	descriptionAlternativeCoverRegex = regexp.MustCompile(`^<i>.*?[Aa]lternat(iv)?e [Cc]over.*?</i>`)

	breakTagRegex = regexp.MustCompile(`<br\s*/?>`)
)

type BookOverview struct {
	Id        string `xml:"id"`
	FullTitle string `xml:"title"`
	Author    string `xml:"author>name"`
}

// Title is the full title with any subtitle and series removed.
func (o BookOverview) Title() string {
	return extractTitle(o.FullTitle)
}

// Subtitle is the subtitle part of the full title with any series removed.
func (o BookOverview) Subtitle() string {
	return extractSubtitle(o.FullTitle)
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

// Title is the full title with any subtitle and series removed.
func (w Work) Title() string {
	return extractTitle(w.FullTitle)
}

// Subtitle is the subtitle part of the full title with any series removed.
func (w Work) Subtitle() string {
	return extractSubtitle(w.FullTitle)
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
	PublicationYear  int     `xml:"publication_year"`
	PublicationMonth int     `xml:"publication_month"`
	PublicationDay   int     `xml:"publication_day"`
	Publisher        string  `xml:"publisher"`
	CountryCode      string  `xml:"country_code"`
	Language         string  `xml:"language_code"`
}

// Title is the full title with any subtitle and series removed.
func (e Edition) Title() string {
	return extractTitle(e.FullTitle)
}

// Subtitle is the subtitle part of the full title with any series removed.
func (e Edition) Subtitle() string {
	return extractSubtitle(e.FullTitle)
}

func (e *Edition) Sanitise() {
	// Description is html and can contain preamble about alternative covers.
	// Break tags need to be specially handled to add new lines as html2text does
	// not convert them to new lines properly
	e.Description = descriptionAlternativeCoverRegex.ReplaceAllString(e.Description, "")
	e.Description = breakTagRegex.ReplaceAllString(e.Description, "\n")
	e.Description = html2text.HTML2TextWithOptions(e.Description, html2text.WithUnixLineBreaks())
	e.Description = strings.TrimSpace(e.Description)

	// Get original cover image by cleaning the url
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
		titles = append(titles, book.Title())
	}
	return titles
}

// extractTitle extracts the title from the full title with any subtitle and series removed.
func extractTitle(fullTitle string) string {
	titleParts := strings.Split(fullTitle, ":")

	title := titleParts[0]
	title = strings.TrimSpace(title)
	title = titleSeriesRegex.ReplaceAllString(title, "")
	title = strings.TrimSpace(title)

	return title
}

// extractTitle extracts the subtitle part of the full title with any series removed.
func extractSubtitle(fullTitle string) string {
	colonIdx := strings.Index(fullTitle, ":")
	if colonIdx == -1 {
		return ""
	}

	subtitle := fullTitle[colonIdx+1:]
	subtitle = strings.TrimSpace(subtitle)
	subtitle = titleSeriesRegex.ReplaceAllString(subtitle, "")
	subtitle = strings.TrimSpace(subtitle)

	return subtitle
}
