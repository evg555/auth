package user

import (
	"context"
	"github.com/evg555/auth/internal/model"
)

func (s *srv) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64

	id, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return id, err
	}

	return id, nil
}
