package user

import (
	"github.com/evg555/auth/internal/repository"
	"github.com/evg555/auth/internal/service"
)

type srv struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) service.UserService {
	return &srv{userRepository: userRepository}
}
