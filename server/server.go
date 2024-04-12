package server

import (
	"context"

	"github.com/ahobsonsayers/abs-goodreads/goodreads"
	"github.com/ahobsonsayers/abs-goodreads/utils.go"
)

type server struct{}

func NewServer() StrictServerInterface { return &server{} }

func (s *server) Search(ctx context.Context, request SearchRequestObject) (SearchResponseObject, error) {
	// Search book
	goodreadsBooks, err := goodreads.DefaultGoodreadsClient.SearchBook(ctx, request.Params.Query, request.Params.Author)
	if err != nil {
		return Search500JSONResponse{Error: utils.ToPointer(err.Error())}, nil
	}

	books := make([]BookMetadata, 0, len(goodreadsBooks))
	for _, goodreadsBook := range goodreadsBooks {
		book := GoodreadsBookToAudioBookShelfBook(goodreadsBook)
		books = append(books, book)
	}

	return Search200JSONResponse{Matches: &books}, nil
}
