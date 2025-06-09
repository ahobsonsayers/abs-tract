package server

import (
	"context"

	"github.com/samber/lo"
)

type server struct{}

func NewServer() StrictServerInterface { return &server{} }

func (*server) SearchGoodreads(
	ctx context.Context,
	request SearchGoodreadsRequestObject,
) (SearchGoodreadsResponseObject, error) {
	books, err := searchGoodreadsBooks(ctx, request.Params.Query, request.Params.Author)
	if err != nil {
		return SearchGoodreads500JSONResponse{N500JSONResponse{Error: lo.ToPtr(err.Error())}}, nil
	}

	return SearchGoodreads200JSONResponse{N200JSONResponse{Matches: &books}}, nil
}

func (*server) SearchKindle(
	ctx context.Context,
	request SearchKindleRequestObject,
) (SearchKindleResponseObject, error) {
	books, err := searchKindleBooks(ctx, request.Region, request.Params.Query, request.Params.Author)
	if err != nil {
		return SearchKindle500JSONResponse{N500JSONResponse{Error: lo.ToPtr(err.Error())}}, nil
	}

	return SearchKindle200JSONResponse{N200JSONResponse{Matches: &books}}, nil
}

func (*server) SearchBookbeat(
	ctx context.Context,
	request SearchBookbeatRequestObject,
) (SearchBookbeatResponseObject, error) {
	books, err := searchBookbeatBooks(ctx,
		request.Market,
		request.Format,
		request.Languages,
		request.Params.Query,
		request.Params.Author)
	if err != nil {
		return SearchBookbeat500JSONResponse{N500JSONResponse{Error: lo.ToPtr(err.Error())}}, nil
	}

	return SearchBookbeat200JSONResponse{N200JSONResponse{Matches: &books}}, nil
}
