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
	// Shelves is a struct matching the goodread response xml
	var shelves struct {
		Shelf []struct {
			Name string `xml:"name,attr"`
		} `xml:"shelf"`
	}
	err := d.DecodeElement(&shelves, &start)
	if err != nil {
		return err
	}

	genres := make(Genres, 0, 5)
	seenGenreShelves := mapset.NewSetWithSize[string](5)
	for _, shelf := range shelves.Shelf {
		// Make shelf name singular for easier comparison
		shelfName := inflection.Singular(shelf.Name)

		// Skip non genre shelves and already seen genre shelves
		if !genreShelves.Contains(shelfName) || seenGenreShelves.Contains(shelfName) {
			continue
		}

		// Get genre from shelf name, making it human readable
		genre := strings.ReplaceAll(shelfName, "-", " ")
		genreWords := strings.Fields(genre)
		for i, word := range genreWords {
			genreWords[i] = cases.Title(language.Und).String(word) // Capitalize each word
		}
		genre = strings.Join(genreWords, " ")

		genres = append(genres, genre)
		seenGenreShelves.Add(shelfName)
		if len(genres) == 5 {
			break
		}
	}
	*g = genres

	return nil
}
