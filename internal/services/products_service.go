package services

import (
	"context"
	"fmt"
	sqlc "simplestore_go/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) GetProducts(ctx context.Context, args sqlc.ListProductsParams) ([]sqlc.ListProductsRow, error) {
	Products, err := s.queries.ListProducts(ctx, args)
	fmt.Println(err)
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

func (s *Service) CreateProduct(ctx context.Context, args sqlc.CreateProductParams) (sqlc.Product, error) {
	return s.queries.CreateProduct(ctx, args)
}
func (s *Service) UpdateProduct(ctx context.Context, args sqlc.UpdateProductParams) (sqlc.Product, error) {
	return s.queries.UpdateProduct(ctx, args)
}
func (s *Service) DeleteProduct(ctx context.Context, id pgtype.UUID) error {
	return s.queries.DeleteProduct(ctx, id)
}
