package db

import (
	"context"
	"embed"
	"fmt"
	"irule-api/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations
var embedMigrations embed.FS

func NewPgPool(cfg *config.Config) (*pgxpool.Pool, error) {
	user := cfg.DB_USER
	password := cfg.DB_PASSWORD
	host := cfg.DB_HOST
	port := cfg.DB_PORT
	dbname := cfg.DB_NAME

	connString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", user, password, host, port, dbname)
	pool, err := pgxpool.New(context.Background(), connString)

	if err != nil {
		zap.S().Fatalw("failed to create pgx pool", err)
		return nil, fmt.Errorf("failed to create pgx pool: %s", err)
	}

	var test string

	err = pool.QueryRow(context.Background(), "SELECT 'Hello, world!'").Scan(&test)
	if err != nil {
		zap.S().Fatalw("failed to connect to database", err)
		return nil, fmt.Errorf("failed to connect to database: %s", err)
	}
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		zap.S().Fatalw("failed to set goose dialect", err)
		return nil, fmt.Errorf("failed to set goose dialect: %s", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(db, "migrations"); err != nil {
		zap.S().Fatalw("failed to run goose up", err)
		return nil, fmt.Errorf("failed to run goose up: %s", err)
	}
	if err := db.Close(); err != nil {
		zap.S().Fatalw("failed to close db", err)
		return nil, fmt.Errorf("failed to close db: %s", err)
	}

	zap.S().Info("database migrations ran successfully")
	return pool, nil
}
