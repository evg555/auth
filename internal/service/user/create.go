package user

import (
	"context"
	"github.com/evg555/auth/internal/model"
)

func (s *srv) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64

	err := s.txManager.ReadComitted(ctx, func(ctx context.Context) error {
		var errTx error

		id, errTx = s.userRepository.Create(ctx, user)
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.Log(ctx, MethodCreate, user)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
