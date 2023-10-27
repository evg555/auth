package api

import (
	"context"

	"github.com/evg555/auth/internal/converter"
	proto "github.com/evg555/auth/pkg/user_v1"
)

func (s *Server) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	user, err := s.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToGetResponseFromService(user), nil
}
