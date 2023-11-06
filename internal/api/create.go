package api

import (
	"context"

	"github.com/evg555/auth/internal/converter"
	proto "github.com/evg555/auth/pkg/user_v1"
)

func (s *Server) Create(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	id, err := s.userService.Create(ctx, converter.ToUserFromCreateRequest(req))
	if err != nil {
		return nil, err
	}

	return &proto.CreateResponse{Id: id}, nil
}
