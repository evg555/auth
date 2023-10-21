package api

import (
	"context"
	proto "github.com/evg555/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (s *Server) Delete(ctx context.Context, req *proto.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("deleting user with id: %d\n", req.GetId())

	err := s.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
