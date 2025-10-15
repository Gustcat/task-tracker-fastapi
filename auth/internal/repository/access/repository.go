package access

import (
	"auth/internal/client/db"
	"auth/internal/logger"
	"auth/internal/repository"
	"context"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

const (
	tableName = "accesses"

	endpointColumn = "endpoint"
	roleColumn     = "role"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{db: db}
}

func (r *repo) Check(ctx context.Context, role int32, endpoint string) error {
	builder := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.And{
			sq.Eq{endpointColumn: endpoint},
			sq.Eq{roleColumn: role},
		}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error("Sql query is not generated", zap.String("endpoint", endpoint), zap.Int32("role", role))
		return err
	}

	q := db.Query{
		Name:     "auth_repository.Login",
		QueryRaw: query,
	}

	ct, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		logger.Error("Query failed", zap.String("endpoint", endpoint), zap.Int32("role", role))
		return err
	}

	if ct.RowsAffected() == 0 {
		logger.Info("user has no access to the endpoint", zap.String("endpoint", endpoint), zap.Int32("role", role))
		return err
	}

	return nil
}
