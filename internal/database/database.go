package database

import (
	"context"
	"fmt"

	"github.com/davidcm146/assets-management-be.git/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(ctx context.Context, dbCfg config.DBConfig) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName,
	)

	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create db pool: %w", err)
	}

	return pool, nil
}
