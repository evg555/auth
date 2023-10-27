package api

import (
	"context"
	"github.com/evg555/auth/internal/converter"
	proto "github.com/evg555/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Update(ctx context.Context, req *proto.UpdateRequest) (*emptypb.Empty, error) {
	err := s.userService.Update(ctx, converter.ToUserFromUpdateRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
