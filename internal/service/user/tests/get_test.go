package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

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

func TestGet(t *testing.T) {
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

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, false, false, 10)
		role     = gofakeit.IntRange(0, 1)

		repoErr = errors.New("repository error")
		method  = "get"

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
		want             *model.User
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
			want: user,
			err:  nil,
			userRepoMockFunc: func(ctx context.Context, mc *minimock.Controller) repository.UserRepository {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
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
			name: "repo error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repoErr,
			userRepoMockFunc: func(ctx context.Context, mc *minimock.Controller) repository.UserRepository {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
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
			name: "log error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repoErr,
			userRepoMockFunc: func(ctx context.Context, mc *minimock.Controller) repository.UserRepository {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				mock.LogMock.Expect(ctx, method, user).Return(repoErr)
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

			newUser, err := serv.Get(tt.args.ctx, tt.args.id)
			err = errors.Unwrap(errors.Unwrap(err))

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newUser)
		})
	}
}
