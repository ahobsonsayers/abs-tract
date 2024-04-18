package server

import (
	"strconv"

	"github.com/ahobsonsayers/abs-goodreads/goodreads"
	"github.com/ahobsonsayers/abs-goodreads/utils"
)

func GoodreadsBookToAudioBookShelfBook(goodreadsBook goodreads.Book) BookMetadata {
	var authorName *string
	if len(goodreadsBook.Authors) != 0 {
		authorName = &goodreadsBook.Authors[0].Name
	}

	series := make([]SeriesMetadata, 0, len(goodreadsBook.Series))
	for _, goodreadsSeriesSingle := range goodreadsBook.Series {
		seriesSingle := SeriesMetadata{
			Series:   goodreadsSeriesSingle.Series.Title,
			Sequence: goodreadsSeriesSingle.BookPosition,
		}
		series = append(series, seriesSingle)
	}

	return BookMetadata{
		// Work Fields
		Title:         goodreadsBook.Work.Title,
		Author:        authorName,
		PublishedYear: utils.ToPointer(strconv.Itoa(goodreadsBook.Work.PublicationYear)),
		// Edition Fields
		Isbn:        goodreadsBook.BestEdition.ISBN,
		Cover:       utils.ToPointer(goodreadsBook.BestEdition.ImageURL),
		Description: &goodreadsBook.BestEdition.Description,
		Publisher:   &goodreadsBook.BestEdition.Publisher,
		Language:    &goodreadsBook.BestEdition.LanguageCode,
		// Other fields
		Series: &series,
		Genres: utils.ToPointer([]string(goodreadsBook.Genres)),
	}
}
