package user

import (
	"auth/internal/client/db"
	"auth/internal/logger"
	"auth/internal/model"
	"auth/internal/repository"
	"auth/internal/repository/user/converter"
	modelRepo "auth/internal/repository/user/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	tableName = "auth_user"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	passwordColumn  = "password"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.UserInfo, pwd string) (int64, error) {
	hashpwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Can't hash the password", zap.String("password", pwd))
		return 0, err
	}

	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn, passwordColumn).
		Values(info.Name, info.Email, info.Role, hashpwd).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error("Sql query to create user is not generated",
			zap.String("name", info.Name),
			zap.String("email", info.Email),
			zap.Int32("role", info.Role),
			zap.String("password", pwd))
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		logger.Error("Query to create user failed",
			zap.String("name", info.Name),
			zap.String("email", info.Email),
			zap.Int32("role", info.Role),
			zap.String("password", pwd))
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (int64, *model.UserInfo, time.Time, sql.NullTime, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn, passwordColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, nil, time.Time{}, sql.NullTime{Time: time.Time{}, Valid: false},
			fmt.Errorf("building SQL failed: %w", err)
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil, time.Time{}, sql.NullTime{Time: time.Time{}, Valid: false},
			repository.ErrUserNotFound
	}
	if err != nil {
		return 0, nil, time.Time{}, sql.NullTime{Time: time.Time{}, Valid: false},
			fmt.Errorf("executing query failed: %w", err)
	}

	id, userInfo, createdAt, updatedAt := converter.ToUserFromRepo(&user)
	return id, userInfo, createdAt, updatedAt, nil
}

func (r *repo) Update(ctx context.Context, id int64, name string, email string) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{"id": id})

	if name != "" {
		builder = builder.Set(nameColumn, name)
	}

	if email != "" {
		builder = builder.Set(emailColumn, email)

	}

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error("Sql query to update user is not generated",
			zap.Int64("id", id),
			zap.String("name", name),
			zap.String("email", email))
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	ct, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		logger.Error("Query to update user failed",
			zap.Int64("id", id),
			zap.String("name", name),
			zap.String("email", email))
		return err
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("запись с id %d не найдена", id)
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error("Sql query to delete user is not generated", zap.Int64("id", id))
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	ct, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		logger.Error("Query to delete user failed", zap.Int64("id", id))
		return err
	}

	if ct.RowsAffected() == 0 {
		return repository.ErrUserNotFound
	}

	return nil
}

func (r *repo) Login(ctx context.Context, username string) (string, *model.UserToken, error) {
	builder := sq.Select(idColumn, roleColumn, passwordColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: username}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error("Sql query to login user is not generated", zap.String("username", username))
		return "", nil, err
	}

	q := db.Query{
		Name:     "user_repository.Login",
		QueryRaw: query,
	}

	user := &model.UserToken{
		Name: username,
	}
	var hashPassword string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Role, &hashPassword)
	if err != nil {
		logger.Error("Query to login user failed", zap.String("username", username))
		return "", nil, err
	}

	return hashPassword, user, nil
}
