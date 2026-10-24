package services

import (
	sqlc "simplestore_go/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	queries *sqlc.Queries
}

func NewService(db *pgxpool.Pool) *Service {
	queries := sqlc.New(db)
	return &Service{queries: queries}
}
