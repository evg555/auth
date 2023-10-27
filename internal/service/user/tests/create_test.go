package tests

import (
	"context"
	"database/sql"
	"github.com/evg555/auth/internal/client/db"
	"github.com/evg555/auth/internal/client/db/transaction"
	"github.com/evg555/auth/internal/repository"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	dbMock "github.com/evg555/auth/internal/client/db/mocks"
	"github.com/evg555/auth/internal/model"
	repoMock "github.com/evg555/auth/internal/repository/mocks"
	userServ "github.com/evg555/auth/internal/service/user"
)

func TestCreate(t *testing.T) {
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

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, false, false, 10)
		role     = gofakeit.IntRange(0, 1)

		repoErr = errors.New("repository error")
		method  = "create"
		txKey   = "tx"

		opts = pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

		user = &model.User{
			ID: 0,
			Name: sql.NullString{
				String: name,
				Valid:  true,
			},
			Password: password,
			Email: sql.NullString{
				String: email,
				Valid:  true,
			},
			Role:      int32(role),
			CreatedAt: time.Time{},
			UpdatedAt: sql.NullTime{},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name             string
		args             args
		want             int64
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
			want: id,
			err:  nil,
			userRepoMockFunc: func(ctx context.Context, mc *minimock.Controller) repository.UserRepository {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(id, nil)
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
			want: 0,
			err:  repoErr,
			userRepoMockFunc: func(ctx context.Context, mc *minimock.Controller) repository.UserRepository {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(0, repoErr)
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

			newID, err := serv.Create(tt.args.ctx, tt.args.user)
			err = errors.Unwrap(errors.Unwrap(err))

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
