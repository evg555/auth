package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/evg555/platform-common/pkg/db"
	"github.com/evg555/platform-common/pkg/db/pg"
	"github.com/evg555/platform-common/pkg/db/transaction"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/evg555/auth/internal/model"
	"github.com/evg555/auth/internal/repository"
	repoMock "github.com/evg555/auth/internal/repository/mocks"
	dbMock "github.com/evg555/auth/internal/service/mocks"
	userServ "github.com/evg555/auth/internal/service/user"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type userRepoMockFunc func(ctx context.Context, mc *minimock.Controller) repository.UserRepository
	type txManagerFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx  context.Context
		user *model.User
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()

		repoErr = errors.New("repository error")
		method  = "update"

		opts = pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

		user = &model.User{
			ID: id,
			Name: sql.NullString{
				String: name,
				Valid:  true,
			},
			Email: sql.NullString{
				String: email,
				Valid:  true,
			},
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
				ctx:  ctx,
				user: user,
			},
			err: nil,
			userRepoMockFunc: func(ctx context.Context, mc *minimock.Controller) repository.UserRepository {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, user).Return(nil)
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
				ctx:  ctx,
				user: user,
			},
			err: repoErr,
			userRepoMockFunc: func(ctx context.Context, mc *minimock.Controller) repository.UserRepository {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, user).Return(repoErr)
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
			ctxNew := context.WithValue(ctx, pg.TxKey, TxMock{})
			userRepoMock := tt.userRepoMockFunc(ctxNew, mc)
			serv := userServ.NewService(userRepoMock, txManagerMock)

			err := serv.Update(tt.args.ctx, tt.args.user)
			err = errors.Unwrap(errors.Unwrap(err))

			require.Equal(t, tt.err, err)
		})
	}
}
