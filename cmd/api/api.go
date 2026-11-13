package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"simplestore_go/internal/routes"
	"simplestore_go/internal/search"
	"simplestore_go/internal/services"
)

func RegisterRoutes(mux *chi.Mux, service services.Service) {
	routes.RegisterUsersRoutes(mux, service)
	routes.RegisterCategoriesRoutes(mux, service)
	routes.RegisterProductsRoutes(mux, service)
	routes.RegisterFilesRoutes(mux, service)
	osClient := search.NewClient()

	mux.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		resp, err := search.SearchProducts(r.Context(), osClient, q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		json.NewEncoder(w).Encode(resp.Hits.Hits)
	})
}

func Mount(db *pgxpool.Pool, mux *chi.Mux) {
	service := services.NewService(db)
	RegisterRoutes(mux, *service)
}

func DBInit() *pgxpool.Pool {
	ctx := context.TODO()
	dbSource := "postgres://root:" + os.Getenv(
		"DATABASE_PASSWORD",
	) + "@" + os.Getenv(
		"DATABASE_HOST",
	) + ":5432/" + os.Getenv("DATABASE_NAME") + "?sslmode=disable"
	db, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Successfully Connect to Database")
	return db
}
