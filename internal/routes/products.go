package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	sqlc "simplestore_go/internal/db"
	"simplestore_go/internal/services"
	"simplestore_go/utils"
)

func RegisterProductsRoutes(mux *chi.Mux, service services.Service) {
	mux.With(utils.AdminOrReadOnlyMiddleware).Route("/products", func(r chi.Router) {
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			id := chi.URLParam(r, "id")
			var uuidId pgtype.UUID
			err := uuidId.Scan(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			product, err := service.GetProduct(r.Context(), uuidId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = json.NewEncoder(w).Encode(product)
			w.WriteHeader(http.StatusCreated)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		})

		r.Get("/by_slug/{slug}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			rawSlug := r.PathValue("slug")

			fmt.Println(rawSlug)
			// chi may return an escaped value when RawPath is set.
			slug, err := url.PathUnescape(rawSlug)
			fmt.Println(slug)

			if err != nil {
				http.Error(w, "invalid slug encoding", http.StatusBadRequest)
				return
			}

			product, err := service.GetProductBySlug(
				r.Context(),
				pgtype.Text{String: slug, Valid: true},
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := json.NewEncoder(w).Encode(product); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		})

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var products []sqlc.ListProductsRow
			products, err := service.GetProducts(
				r.Context(),
				sqlc.ListProductsParams{Limit: 100, Offset: 0},
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = json.NewEncoder(w).Encode(products)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
		// r.Get("/search/{keyword}", func(w http.ResponseWriter, r *http.Request) {
		// 	w.Header().Set("Content-Type", "application/json")
		// 	var products []sqlc.SearchProductsRow
		//
		// 	keyword := chi.URLParam(r, "keyword")
		// 	products, err := service.SearchProduct(r.Context(), keyword)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 		return
		// 	}
		// 	err = json.NewEncoder(w).Encode(products)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 		return
		// 	}
		// })
		//
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var product services.CereateProductDto
			err := json.NewDecoder(r.Body).Decode(&product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, err = service.CreateProduct(r.Context(), product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
		})
		r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			id := chi.URLParam(r, "id")
			var product services.UpdateProductDto
			var uuidId pgtype.UUID
			err := uuidId.Scan(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			product.ID = uuidId
			err = json.NewDecoder(r.Body).Decode(&product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, err = service.UpdateProduct(r.Context(), product)
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
			err = service.DeleteProduct(r.Context(), uuidId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})
	})
}
