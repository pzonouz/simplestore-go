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
	args sqlc.UpdateProductParams,
) (sqlc.Product, error) {
	// args.Description.String = utils.NormalizeSlug(args.Description.String)
	return s.queries.UpdateProduct(ctx, args)
}

func (s *Service) DeleteProduct(ctx context.Context, id pgtype.UUID) error {
	return s.queries.DeleteProduct(ctx, id)
}
