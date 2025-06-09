package server

import (
	"context"
	"strconv"

	"github.com/ahobsonsayers/abs-tract/bookbeat"
	"github.com/ahobsonsayers/abs-tract/goodreads"
	"github.com/ahobsonsayers/abs-tract/kindle"
	"github.com/samber/lo"
)

func searchGoodreadsBooks(ctx context.Context, title string, author *string) ([]BookMetadata, error) {
	goodreadsBooks, err := goodreads.DefaultClient.SearchBooks(ctx, title, author)
	if err != nil {
		return nil, err
	}

	// Limit number of books to 20
	if len(goodreadsBooks) > 20 {
		goodreadsBooks = goodreadsBooks[:20]
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
	title string,
	author *string,
) ([]BookMetadata, error) {
	kindleClient, err := kindle.NewClient(lo.ToPtr(string(countryCode)))
	if err != nil {
		return nil, err
	}

	kindleBooks, err := kindleClient.Search(ctx, title, author)
	if err != nil {
		return nil, err
	}

	// Limit number of books to 20
	if len(kindleBooks) > 20 {
		kindleBooks = kindleBooks[:20]
	}

	books := make([]BookMetadata, 0, len(kindleBooks))
	for _, kindleBook := range kindleBooks {
		book := kindleBookToBookMetadata(kindleBook)
		books = append(books, book)
	}

	return books, nil
}

func searchBookbeatBooks(
	ctx context.Context,
	marketStr BookbeatMarket,
	formatStr BookbeatFormat,
	languagesStr BookbeatLanguageCodes,
	title string,
	author *string,
) ([]BookMetadata, error) {

	market := string(marketStr)
	format := string(formatStr)
	languages := string(languagesStr)
	bookbeatClient, err := bookbeat.NewClient(market, format, languages)
	if err != nil {
		return nil, err
	}

	bookbeatBooks, err := bookbeatClient.Search(ctx, title, author)
	if err != nil {
		return nil, err
	}

	books := make([]BookMetadata, 0, len(bookbeatBooks))
	for _, bookbeatBook := range bookbeatBooks {
		book := bookbeatBookToBookMetadata(bookbeatBook)
		books = append(books, book)
	}

	return books, nil
}

func goodreadsBookToBookMetadata(goodreadsBook goodreads.Book) BookMetadata {
	var subtitle *string
	if goodreadsBook.BestEdition.Subtitle() != "" {
		subtitle = lo.ToPtr(goodreadsBook.BestEdition.Subtitle())
	}

	var author *string
	if len(goodreadsBook.Authors) != 0 {
		author = &goodreadsBook.Authors[0].Name
	}

	var publicationYear *string
	if goodreadsBook.Work.PublicationYear != 0 {
		publicationYear = lo.ToPtr(strconv.Itoa(goodreadsBook.Work.PublicationYear))
	} else if goodreadsBook.BestEdition.PublicationYear != 0 {
		publicationYear = lo.ToPtr(strconv.Itoa(goodreadsBook.BestEdition.PublicationYear))
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
		Title:         goodreadsBook.BestEdition.Title(),
		Subtitle:      subtitle,
		Author:        author,
		PublishedYear: publicationYear,
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

func bookbeatBookToBookMetadata(bookbeatBook bookbeat.Book) BookMetadata {
	var publishedYear *string
	publishedYear = lo.ToPtr(strconv.Itoa(bookbeatBook.PublishedYear))

	var series []SeriesMetadata
	if bookbeatBook.Series != nil {
		series = make([]SeriesMetadata, 0, 1)
		series = append(series, SeriesMetadata{
			Series:   bookbeatBook.Series.Series,
			Sequence: lo.ToPtr(strconv.Itoa(bookbeatBook.Series.Sequence)),
		})
	}

	return BookMetadata{
		Author:        &bookbeatBook.Authors,
		Cover:         &bookbeatBook.Cover,
		Description:   &bookbeatBook.Description,
		Duration:      &bookbeatBook.Duration,
		Genres:        &bookbeatBook.Genres,
		Isbn:          &bookbeatBook.ISBN,
		Language:      &bookbeatBook.Language,
		Narrator:      &bookbeatBook.Narrators,
		PublishedYear: publishedYear,
		Publisher:     &bookbeatBook.Publisher,
		Series:        &series,
		Tags:          &bookbeatBook.Tags,
		Title:         bookbeatBook.Title,
		Subtitle:      &bookbeatBook.Subtitle,
	}
}
