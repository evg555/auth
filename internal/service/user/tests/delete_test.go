package tests

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/evg555/auth/internal/client/db"
	"github.com/evg555/auth/internal/client/db/transaction"
	"github.com/evg555/auth/internal/repository"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"

	dbMock "github.com/evg555/auth/internal/client/db/mocks"
	"github.com/evg555/auth/internal/model"
	repoMock "github.com/evg555/auth/internal/repository/mocks"
	userServ "github.com/evg555/auth/internal/service/user"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type userRepoMockFunc func(ctx context.Context, mc *minimock.Controller) repository.UserRepository
	type txManagerFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoErr = errors.New("repository error")
		method  = "delete"
		txKey   = "tx"

		opts = pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

		user = &model.User{
			ID: id,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name             string
		args             args
		err              error
		userRepoMockFunc userRepoMockFunc
		txManagerFunc    txManagerFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: nil,
			userRepoMockFunc: func(ctx context.Context, mc *minimock.Controller) repository.UserRepository {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				mock.LogMock.Expect(ctx, method, user).Return(nil)
				return mock
			},
			txManagerFunc: func(mc *minimock.Controller) db.TxManager {
				mock := dbMock.NewTransactorMock(mc)
				mock.BeginTxMock.Expect(ctx, opts).Return(TxMock{}, nil)

				txManager := transaction.NewTransactionManager(mock)
				return txManager
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: repoErr,
			userRepoMockFunc: func(ctx context.Context, mc *minimock.Controller) repository.UserRepository {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(repoErr)
				return mock
			},
			txManagerFunc: func(mc *minimock.Controller) db.TxManager {
				mock := dbMock.NewTransactorMock(mc)
				mock.BeginTxMock.Expect(ctx, opts).Return(TxMock{}, nil)

				txManager := transaction.NewTransactionManager(mock)
				return txManager
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			txManagerMock := tt.txManagerFunc(mc)
			ctxNew := context.WithValue(ctx, txKey, TxMock{})
			userRepoMock := tt.userRepoMockFunc(ctxNew, mc)
			serv := userServ.NewService(userRepoMock, txManagerMock)

			err := serv.Delete(tt.args.ctx, tt.args.id)
			err = errors.Unwrap(errors.Unwrap(err))

			require.Equal(t, tt.err, err)
		})
	}
}
