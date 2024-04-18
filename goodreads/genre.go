package goodreads

import (
	"encoding/xml"
	"strings"

	"github.com/jinzhu/inflection"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	mapset "github.com/deckarep/golang-set/v2"
)

func init() {
	// Add genre shelves to uncountable inflections
	inflection.AddUncountable(genreShelves.ToSlice()...)
}

// None exhaustive list of popular genre shelves. More should be added.
// Most of these are obtained from https://www.goodreads.com/genres and made singular.
// Used to allow extraction of genres from goodreads shelves.
// TODO is there a better way to recognise genres?
var genreShelves = mapset.NewSet(
	"adventure",
	"art",
	"biography",
	"business",
	"chick-lit",
	"childrens",
	"christian",
	"classic",
	"comedy",
	"comic",
	"contemporary",
	"cookbook",
	"crime",
	"fantasy",
	"fiction",
	"gay-and-lesbian",
	"graphic-novel",
	"high-fantasy",
	"historical-fiction",
	"historical",
	"history",
	"horror",
	"humor-and-comedy",
	"humour",
	"manga",
	"memoir",
	"music",
	"mystery",
	"non-fiction",
	"nonfiction",
	"paranormal",
	"philosophy",
	"picture-book",
	"poetry",
	"politics",
	"psychology",
	"religion",
	"romance",
	"science-fiction",
	"science",
	"self-help",
	"spirituality",
	"sport",
	"suspense",
	"thriller",
	"travel",
	"young-adult",
)

type Genres []string

func (g *Genres) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// unmarshaller is a struct matching the goodreads response xml
	var unmarshaller struct {
		Shelves []struct { // nolint
			Name string `xml:"name,attr"`
		} `xml:"shelf"`
	}
	err := d.DecodeElement(&unmarshaller, &start)
	if err != nil {
		return err
	}

	// Get shelf names
	shelfNames := make([]string, 0, len(unmarshaller.Shelves))
	for _, shelf := range unmarshaller.Shelves {
		shelfNames = append(shelfNames, shelf.Name)
	}

	// Convert shelf names to genres
	genres := shelvesToGenres(shelfNames)

	// Only use first (up to) three genres
	if len(genres) < 3 {
		*g = genres
	} else {
		*g = genres[:3]
	}

	return nil
}

func shelvesToGenres(shelves []string) []string {
	genres := make(Genres, 0, len(shelves))
	seenGenreShelves := mapset.NewSetWithSize[string](len(shelves))
	for _, shelf := range shelves {
		// Make shelf singular for easier comparison
		shelf := inflection.Singular(shelf)

		// Skip non genre shelves and already seen genre shelves
		if !genreShelves.Contains(shelf) || seenGenreShelves.Contains(shelf) {
			continue
		}

		// Get genre from shelf name, making it human readable
		genre := strings.ReplaceAll(shelf, "-", " ")
		genreWords := strings.Fields(genre)
		for i, word := range genreWords {
			genreWords[i] = cases.Title(language.Und).String(word) // Capitalize each word
		}
		genre = strings.Join(genreWords, " ")

		genres = append(genres, genre)
		seenGenreShelves.Add(shelf)
	}

	return genres
}
