package user

import (
	"context"
	"github.com/evg555/auth/internal/model"
)

func (s *srv) Update(ctx context.Context, user *model.User) error {
	err := s.txManager.ReadComitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Update(ctx, user)
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.Log(ctx, MethodUpdate, user)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
