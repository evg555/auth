package user

import (
	"github.com/evg555/auth/internal/client/db"
	"github.com/evg555/auth/internal/repository"
	"github.com/evg555/auth/internal/service"
)

const (
	MethodCreate = "create"
	MethodUpdate = "update"
	MethodGet    = "get"
	MethodDelete = "delete"
)

type srv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) service.UserService {
	return &srv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
