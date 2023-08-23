package server

import (
	"context"
	"log/slog"

	servicepb "github.com/fillmore-labs/name-service/api/fillmore-labs/name-service/v1alpha1"
	"github.com/fillmore-labs/name-service/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type nameServiceServer struct {
	db database.Database
	servicepb.UnimplementedNameServiceServer
}

func NewServer(db database.Database) servicepb.NameServiceServer {
	return &nameServiceServer{db: db}
}

var _ servicepb.NameServiceServer = &nameServiceServer{}

func (s *nameServiceServer) AddName(
	ctx context.Context,
	req *servicepb.AddNameRequest,
) (*servicepb.AddNameResponse, error) {
	name := database.Name{GivenName: req.GetGivenName(), Surname: req.Surname}
	if err := s.db.AddName(ctx, name); err != nil {
		slog.WarnContext(ctx, "can't add name", "err", err)

		return nil, status.Error(codes.Internal, "can't add name")
	}

	return &servicepb.AddNameResponse{}, nil
}

func (s *nameServiceServer) ListNames(
	_ *servicepb.ListNamesRequest,
	stream servicepb.NameService_ListNamesServer,
) error {
	sendName := func(name database.Name) error {
		return stream.Send(
			&servicepb.ListNamesResponse{
				GivenName: name.GivenName,
				Surname:   name.Surname,
			},
		)
	}

	return s.db.ListNames(stream.Context(), sendName)
}
