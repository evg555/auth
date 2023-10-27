package user

import (
	"context"

	"github.com/evg555/auth/internal/model"
)

func (s *srv) Get(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User

	err := s.txManager.ReadComitted(ctx, func(ctx context.Context) error {
		var errTx error

		user, errTx = s.userRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.Log(ctx, MethodGet, user)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
