package server

import (
	"context"

	"github.com/ahobsonsayers/abs-goodreads/goodreads"
	"github.com/samber/lo"
)

type server struct{}

func NewServer() StrictServerInterface { return &server{} }

func (*server) Search(ctx context.Context, request SearchRequestObject) (SearchResponseObject, error) {
	// Search book
	goodreadsBooks, err := goodreads.DefaultClient.SearchBooks(ctx, request.Params.Query, request.Params.Author)
	if err != nil {
		return Search500JSONResponse{Error: lo.ToPtr(err.Error())}, nil
	}

	books := make([]BookMetadata, 0, len(goodreadsBooks))
	for _, goodreadsBook := range goodreadsBooks {
		book := GoodreadsBookToAudioBookShelfBook(goodreadsBook)
		books = append(books, book)
	}

	return Search200JSONResponse{Matches: &books}, nil
}
