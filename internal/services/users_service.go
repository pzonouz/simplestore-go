package services

import (
	"context"
	"fmt"
	sqlc "simplestore_go/internal/db"
	"time"
	"uuid"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) GetUsers(ctx context.Context) ([]sqlc.User, error) {

	users, err := s.queries.ListUsers(ctx, sqlc.ListUsersParams{})
	if err != nil {
		return []sqlc.User{}, err
	}
	if users == nil {
		return []sqlc.User{}, err
	}
	return users, err
}
func (s *Service) CreateUsers(ctx context.Context, args sqlc.CreateUserParams) (sqlc.User, error) {
	token := (uuid.New()).String()
	now := time.Now()
	expiredTime := now.Add(24 * time.Hour)
	pass, err := argon2id.CreateHash(args.Password, argon2id.DefaultParams)
	if err != nil {
		return sqlc.User{}, err
	}
	args.Password = pass
	args.Token = pgtype.Text{String: token, Valid: true}
	args.ExpiredTime = pgtype.Timestamptz{Time: expiredTime, Valid: true}
	return s.queries.CreateUser(ctx, args)
}
func (s *Service) ActivateUserByToken(ctx context.Context, token pgtype.Text) error {
	rowsCount, err := s.queries.ActivateUserByToken(ctx, token)
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return fmt.Errorf("انجام نشد")
	}
	return nil
}
func (s *Service) GetUserByEmail(ctx context.Context, email pgtype.Text) (sqlc.User, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return sqlc.User{}, err
	}
	if !user.IsActive.Bool {
		return sqlc.User{}, fmt.Errorf("")
	}
	return user, nil
}
