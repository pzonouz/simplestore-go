package services

import (
	"github.com/jackc/pgx/v5/pgxpool"

	sqlc "simplestore_go/internal/db"
)

type Service struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	queries := sqlc.New(db)
	return &Service{queries: queries, pool: db}
}
