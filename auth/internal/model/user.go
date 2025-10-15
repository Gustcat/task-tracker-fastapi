package model

import (
	"database/sql"
	"time"
)

// модели сервисов
type User struct {
	ID        int64        `db:"id"`
	Info      UserInfo     `db:""`
	Password  string       `db:"password"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type UserInfo struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  int32  `db:"role"`
}

type UserToken struct {
	Name string `db:"name"`
	Role int32  `db:"role"`
	ID   int64  `db:"id"`
}
