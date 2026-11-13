package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	sqlc "simplestore_go/internal/db"
)

func (s *Service) GetProducts(
	ctx context.Context,
	args sqlc.ListProductsParams,
) ([]sqlc.ListProductsRow, error) {
	Products, err := s.queries.ListProducts(ctx, args)
	if err != nil {
		return []sqlc.ListProductsRow{}, err
	}
	if Products == nil {
		return []sqlc.ListProductsRow{}, err
	}
	return Products, err
}

func (s *Service) GetProduct(ctx context.Context, id pgtype.UUID) (sqlc.GetProductRow, error) {
	Product, err := s.queries.GetProduct(ctx, id)
	if err != nil {
		return sqlc.GetProductRow{}, err
	}
	return Product, err
}

// func (s *Service) SearchProduct(
// 	ctx context.Context,
// 	keyword string,
// ) ([]sqlc.SearchProductsFuzzyRow, error) {
// 	Products, err := s.queries.SearchProductsFuzzy(ctx, keyword)
// 	if err != nil {
// 		return []sqlc.SearchProductsFuzzyRow{}, err
// 	}
// 	return Products, err
// }

func (s *Service) GetProductBySlug(
	ctx context.Context,
	slug pgtype.Text,
) (sqlc.GetProductBySlugRow, error) {
	Product, err := s.queries.GetProductBySlug(ctx, slug)
	if err != nil {
		return sqlc.GetProductBySlugRow{}, err
	}
	return Product, err
}

type CereateProductDto struct {
	Name                    string           `json:"name"`
	Description             pgtype.Text      `json:"description"`
	Category                pgtype.UUID      `json:"category"`
	Slug                    pgtype.Text      `json:"slug"`
	ImageUrl                pgtype.Text      `json:"imageUrl"`
	Count                   pgtype.Numeric   `json:"count"`
	Price                   pgtype.Numeric   `json:"price"`
	PackingType             sqlc.PackingType `json:"packingType"`
	PackingCoefficientTypes []pgtype.Numeric `json:"packingCoefficientTypes"`
	FileIds                 []pgtype.UUID    `json:"fileIds"`
}
type UpdateProductDto struct {
	ID                      pgtype.UUID      `json:"id"`
	Name                    string           `json:"name"`
	Description             pgtype.Text      `json:"description"`
	Category                pgtype.UUID      `json:"category"`
	Slug                    pgtype.Text      `json:"slug"`
	ImageUrl                pgtype.Text      `json:"imageUrl"`
	Count                   pgtype.Numeric   `json:"count"`
	Price                   pgtype.Numeric   `json:"price"`
	PackingType             sqlc.PackingType `json:"packingType"`
	PackingCoefficientTypes []pgtype.Numeric `json:"packingCoefficientTypes"`
	FileIds                 []pgtype.UUID    `json:"fileIds"`
}

func (s *Service) CreateProduct(
	ctx context.Context,
	args CereateProductDto,
) (sqlc.Product, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)
	if err != nil {
		return sqlc.Product{}, err
	}
	q := s.queries.WithTx(tx)
	product, err := q.CreateProduct(ctx, sqlc.CreateProductParams{
		Name:                    args.Name,
		Description:             args.Description,
		Category:                args.Category,
		Slug:                    args.Slug,
		ImageUrl:                args.ImageUrl,
		Count:                   args.Count,
		Price:                   args.Price,
		PackingType:             args.PackingType,
		PackingCoefficientTypes: args.PackingCoefficientTypes,
	})
	if err != nil {
		return sqlc.Product{}, err
	}
	err = q.InsertIntoProductsFiles(
		ctx,
		sqlc.InsertIntoProductsFilesParams{ProductID: product.ID, Column2: args.FileIds},
	)
	if err != nil {
		return sqlc.Product{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return sqlc.Product{}, err
	}
	return product, nil
}

func (s *Service) UpdateProduct(
	ctx context.Context,
	args UpdateProductDto,
) (sqlc.Product, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return sqlc.Product{}, err
	}
	defer tx.Rollback(ctx)
	q := s.queries.WithTx(tx)
	product, err := q.UpdateProduct(
		ctx,
		sqlc.UpdateProductParams{
			Name:                    args.Name,
			Description:             args.Description,
			ImageUrl:                args.ImageUrl,
			Price:                   args.Price,
			Count:                   args.Count,
			Slug:                    args.Slug,
			PackingType:             args.PackingType,
			PackingCoefficientTypes: args.PackingCoefficientTypes,
			ID:                      args.ID,
			Category:                args.Category,
		},
	)
	if err != nil {
		return sqlc.Product{}, err
	}
	err = q.DeleteProductFiles(ctx, args.ID)
	if err != nil {
		return sqlc.Product{}, err
	}
	err = q.InsertIntoProductsFiles(
		ctx,
		sqlc.InsertIntoProductsFilesParams{ProductID: args.ID, Column2: args.FileIds},
	)
	if err != nil {
		return sqlc.Product{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return sqlc.Product{}, err
	}

	return product, nil
}

func (s *Service) DeleteProduct(ctx context.Context, id pgtype.UUID) error {
	return s.queries.DeleteProduct(ctx, id)
}
