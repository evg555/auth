package user

import (
	"context"
	"github.com/evg555/auth/internal/model"
)

func (s *srv) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
