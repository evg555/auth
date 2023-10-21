package api

import (
	"github.com/evg555/auth/internal/service"
	proto "github.com/evg555/auth/pkg/user_v1"
)

type Server struct {
	proto.UnimplementedUserV1Server
	userService service.UserService
}

func NewServer(userService service.UserService) *Server {
	return &Server{
		userService: userService,
	}
}
