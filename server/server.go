package server

import "context"

type server struct{}

func NewServer() StrictServerInterface { return &server{} }

func (s *server) Search(ctx context.Context, request SearchRequestObject) (SearchResponseObject, error) {

	
	return Search500JSONResponse{}, nil
}
