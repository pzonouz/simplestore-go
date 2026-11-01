package main

import (
	"context"
	"log"
	"os"
	"simplestore_go/internal/routes"
	"simplestore_go/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(mux *chi.Mux, service services.Service) {
	routes.RegisterUsersRoutes(mux, service)
	routes.RegisterCategoriesRoutes(mux, service)
	routes.RegisterProductsRoutes(mux, service)
	routes.RegisterFilesRoutes(mux, service)

}

func Mount(db *pgxpool.Pool, mux *chi.Mux) {
	service := services.NewService(db)
	RegisterRoutes(mux, *service)
}
func DBInit() *pgxpool.Pool {
	ctx := context.TODO()
	dbSource := "postgres://root:" + os.Getenv("DATABASE_PASSWORD") + "@" + os.Getenv("DATABASE_HOST") + ":5432/simple_store_go?sslmode=disable"
	db, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Successfully Connect to Database")
	return db
}
