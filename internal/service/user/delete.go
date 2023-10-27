package user

import (
	"context"

	"github.com/evg555/auth/internal/model"
)

func (s *srv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadComitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.Log(ctx, MethodDelete, &model.User{ID: id})
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
