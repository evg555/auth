package service

import (
	"context"

	"github.com/evg555/platform-common/pkg/db"
	"github.com/jackc/pgx/v4"

	"github.com/evg555/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

type Transactor interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (db.Committer, error)
}
