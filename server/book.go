package server

import (
	"strconv"

	"github.com/ahobsonsayers/abs-goodreads/goodreads"
	"github.com/samber/lo"
)

func GoodreadsBookToAudioBookShelfBook(goodreadsBook goodreads.Book) BookMetadata {
	var authorName *string
	if len(goodreadsBook.Authors) != 0 {
		authorName = &goodreadsBook.Authors[0].Name
	}

	var imageUrl *string
	if goodreadsBook.BestEdition.ImageURL != "" {
		imageUrl = lo.ToPtr(goodreadsBook.BestEdition.ImageURL)
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
		PublishedYear: lo.ToPtr(strconv.Itoa(goodreadsBook.Work.PublicationYear)),
		// Edition Fields
		Isbn:        goodreadsBook.BestEdition.ISBN,
		Cover:       imageUrl,
		Description: &goodreadsBook.BestEdition.Description,
		Publisher:   &goodreadsBook.BestEdition.Publisher,
		Language:    &goodreadsBook.BestEdition.Language,
		// Other fields
		Series: &series,
		Genres: lo.ToPtr([]string(goodreadsBook.Genres)),
	}
}
