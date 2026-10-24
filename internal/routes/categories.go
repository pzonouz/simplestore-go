package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	sqlc "simplestore_go/internal/db"
	"simplestore_go/internal/services"
	"simplestore_go/utils"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func RegisterCategoriesRoutes(mux *chi.Mux, service services.Service) {
	mux.With(utils.AdminOrReadOnlyMiddleware).Route("/categories", func(r chi.Router) {
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			id := chi.URLParam(r, "id")
			var uuidId pgtype.UUID
			err := uuidId.Scan(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			category, err := service.GetCategory(r.Context(), uuidId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = json.NewEncoder(w).Encode(category)
			w.WriteHeader(http.StatusCreated)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var categories []sqlc.ListCategoriesRow
			categories, err := service.GetCategories(r.Context(), sqlc.ListCategoriesParams{Limit: 100, Offset: 0})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = json.NewEncoder(w).Encode(categories)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var category sqlc.CreateCategoryParams
			err := json.NewDecoder(r.Body).Decode(&category)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Print(category)
			_, err = service.CreateCategory(r.Context(), category)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
		})
		r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			id := chi.URLParam(r, "id")
			var category sqlc.UpdateCategoryParams
			var uuidId pgtype.UUID
			err := uuidId.Scan(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			category.ID = uuidId
			err = json.NewDecoder(r.Body).Decode(&category)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, err = service.UpdateCategory(r.Context(), category)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
		})
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			var uuidId pgtype.UUID
			err := uuidId.Scan(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = service.DeleteCategory(r.Context(), uuidId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})

	})
}
