package services

import (
	"context"

	sqlc "simple_store_go/internal/db"
	"simple_store_go/myTypes"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) GetCategories(
	ctx context.Context,
	args sqlc.ListCategoriesParams,
) ([]sqlc.ListCategoriesRow, error) {
	Categories, err := s.queries.ListCategories(ctx, args)
	if err != nil {
		return []sqlc.ListCategoriesRow{}, err
	}
	if Categories == nil {
		return []sqlc.ListCategoriesRow{}, err
	}
	return Categories, err
}

func (s *Service) GetSubCategoriesBySlug(
	ctx context.Context,
	slug pgtype.Text,
) ([]sqlc.ListSubCategoriesBySlugRow, error) {
	Categories, err := s.queries.ListSubCategoriesBySlug(ctx, slug)
	if err != nil {
		return []sqlc.ListSubCategoriesBySlugRow{}, err
	}
	if Categories == nil {
		return []sqlc.ListSubCategoriesBySlugRow{}, err
	}
	return Categories, err
}

func (s *Service) GetCategory(ctx context.Context, id pgtype.UUID) (sqlc.GetCategoryRow, error) {
	Category, err := s.queries.GetCategory(ctx, id)
	if err != nil {
		return sqlc.GetCategoryRow{}, err
	}
	return Category, err
}

func (s *Service) CreateCategory(
	ctx context.Context,
	args sqlc.CreateCategoryParams,
) (sqlc.Category, error) {
	return s.queries.CreateCategory(ctx, args)
}

func (s *Service) UpdateCategory(
	ctx context.Context,
	args sqlc.UpdateCategoryParams,
) (sqlc.Category, error) {
	return s.queries.UpdateCategory(ctx, args)
}

func (s *Service) DeleteCategory(ctx context.Context, id pgtype.UUID) error {
	return s.queries.DeleteCategory(ctx, id)
}

func (s *Service) NoChildrenCategories(ctx context.Context) ([]sqlc.Category, error) {
	return s.queries.NoChildrenCategories(ctx)
}

func (s *Service) ListCategoriesTree(
	ctx context.Context,
) ([]myTypes.ListCategoriesTreeRowDTO, error) {
	return s.queries.ListCategoriesTree(ctx)
}
