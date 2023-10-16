package user

import (
	"context"
	"github.com/evg555/auth/internal/model"
)

func (s *srv) Update(ctx context.Context, user *model.User) error {
	err := s.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
