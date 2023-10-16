package user

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/evg555/auth/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"

	"github.com/evg555/auth/internal/repository"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	passwordColumn  = "password"
	emailColumn     = "email"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) repository.UserRepository {
	return &repo{
		conn: conn,
	}
}

func (r *repo) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64

	builderInsert := sq.Insert(tableName).
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

	res := r.conn.QueryRow(ctx, query, args...)

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
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("failed to select from database: %v", err)
		return nil, err
	}

	err = r.conn.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Printf("failed to select from database: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *repo) Update(ctx context.Context, user *model.User) error {
	builderUpdate := sq.Update(tableName).PlaceholderFormat(sq.Dollar)

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

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update database: %v", err)
		return err
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("failed to delete from database: %v", err)
		return err
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete from database: %v", err)
		return err
	}

	return nil
}
