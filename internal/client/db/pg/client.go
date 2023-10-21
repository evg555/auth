package pg

import (
	"context"
	"github.com/evg555/auth/internal/client/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type pgClient struct {
	MasterDB db.DB
}

func New(ctx context.Context, dsn string) (db.Client, error) {
	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &pgClient{MasterDB: NewDB(conn)}, nil
}

func (p *pgClient) DB() db.DB {
	return p.MasterDB
}

func (p *pgClient) Close() error {
	if p.MasterDB != nil {
		p.MasterDB.Close()
	}

	return nil
}
