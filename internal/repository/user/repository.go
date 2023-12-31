package user

import (
	"context"
	"encoding/json"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/evg555/platform-common/pkg/db"

	"github.com/evg555/auth/internal/model"
	"github.com/evg555/auth/internal/repository"
)

const (
	tableNameUsers = "users"

	idColumn        = "id"
	nameColumn      = "name"
	passwordColumn  = "password"
	emailColumn     = "email"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"

	tableNameLogs = "logs"

	methodColumn = "method"
	dataColumn   = "data"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64

	builderInsert := sq.Insert(tableNameUsers).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(
			user.Name.String,
			user.Email.String,
			user.Password,
			user.Role,
		).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to insert to database: %v", err)
		return id, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	res := r.db.DB().QueryRawContext(ctx, q, args...)

	err = res.Scan(&id)
	if err != nil {
		log.Printf("failed to insert to database: %v", err)
		return id, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	var user model.User

	builderSelect := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableNameUsers).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("failed to select from database: %v", err)
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		log.Printf("failed to select from database: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *repo) Update(ctx context.Context, user *model.User) error {
	builderUpdate := sq.Update(tableNameUsers).PlaceholderFormat(sq.Dollar)

	if user.Name.Valid {
		builderUpdate = builderUpdate.Set(nameColumn, user.Name.String)
	}

	if user.Email.Valid {
		builderUpdate = builderUpdate.Set(emailColumn, user.Email.String)
	}

	builderUpdate = builderUpdate.Where(sq.Eq{idColumn: user.ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to update database: %v", err)
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to update database: %v", err)
		return err
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableNameUsers).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("failed to delete from database: %v", err)
		return err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to delete from database: %v", err)
		return err
	}

	return nil
}

func (r *repo) Log(ctx context.Context, method string, user *model.User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	builderInsert := sq.Insert(tableNameLogs).
		PlaceholderFormat(sq.Dollar).
		Columns(methodColumn, dataColumn).
		Values(
			method,
			string(userJson),
		)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to insert to database: %v", err)
		return err
	}

	q := db.Query{
		Name:     "user_repository.Log",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to update database: %v", err)
		return err
	}

	return nil
}
