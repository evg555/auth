package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Name      sql.NullString
	Password  string
	Email     sql.NullString
	Role      int32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
