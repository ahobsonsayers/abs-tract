package server

import (
	"context"
	"strconv"

	"github.com/ahobsonsayers/abs-tract/goodreads"
	"github.com/ahobsonsayers/abs-tract/kindle"
	"github.com/samber/lo"
)

func searchGoodreadsBooks(ctx context.Context, bookTitle string, bookAuthor *string) ([]BookMetadata, error) {
	goodreadsBooks, err := goodreads.DefaultClient.SearchBooks(ctx, bookTitle, bookAuthor)
	if err != nil {
		return nil, err
	}

	books := make([]BookMetadata, 0, len(goodreadsBooks))
	for _, goodreadsBook := range goodreadsBooks {
		book := goodreadsBookToBookMetadata(goodreadsBook)
		books = append(books, book)
	}

	return books, nil
}

func searchKindleBooks(
	ctx context.Context,
	countryCode SearchKindleParamsRegion,
	bookTitle string,
	bookAuthor *string,
) ([]BookMetadata, error) {
	kindleClient, err := kindle.NewClient(nil, lo.ToPtr(string(countryCode)))
	if err != nil {
		return nil, err
	}

	kindleBooks, err := kindleClient.Search(ctx, bookTitle, bookAuthor)
	if err != nil {
		return nil, err
	}

	books := make([]BookMetadata, 0, len(kindleBooks))
	for _, kindleBook := range kindleBooks {
		book := kindleBookToBookMetadata(kindleBook)
		books = append(books, book)
	}

	return books, nil
}

func goodreadsBookToBookMetadata(goodreadsBook goodreads.Book) BookMetadata {
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

func kindleBookToBookMetadata(kindleBook kindle.Book) BookMetadata {
	var publishedYear *string
	if kindleBook.PublishDate != nil {
		publishedYear = lo.ToPtr(strconv.Itoa(kindleBook.PublishDate.Year()))
	}

	return BookMetadata{
		Asin:          &kindleBook.ASIN,
		Title:         kindleBook.Title,
		Author:        &kindleBook.Author,
		Cover:         &kindleBook.Cover,
		PublishedYear: publishedYear,
	}
}
