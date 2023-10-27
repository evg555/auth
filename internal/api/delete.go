package api

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	proto "github.com/evg555/auth/pkg/user_v1"
)

func (s *Server) Delete(ctx context.Context, req *proto.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("deleting user with id: %d\n", req.GetId())

	err := s.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
