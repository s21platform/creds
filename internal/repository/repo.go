package repository

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/s21platform/creds/internal/config"
	"github.com/s21platform/creds/internal/model"
)

type Repository struct {
	db *sqlx.DB
}

func New(cfg *config.Config) *Repository {
	connectCmd := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	db, err := sqlx.Connect("postgres", connectCmd)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &Repository{db}
}

func (r *Repository) GetToken(ctx context.Context, tokenUUID string) (string, error) {
	query, args, err := sq.Select("token").
		From("tokens").
		Where(sq.Eq{"token": tokenUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("failed to build query: %w", err)
	}
	var token string
	err = r.db.GetContext(ctx, &token, query, args...)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	return token, nil
}

func (r *Repository) GetCreds(ctx context.Context, service string) (model.CredsData, error) {
	query, args, err := sq.Select("data").
		From("credentials").
		Where(sq.Eq{"service": service}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return model.CredsData{}, fmt.Errorf("failed to build query: %w", err)
	}
	var creds model.CredsData
	err = r.db.GetContext(ctx, &creds, query, args...)
	if err != nil {
		return model.CredsData{}, fmt.Errorf("failed to get creds: %w", err)
	}
	return creds, nil
}
