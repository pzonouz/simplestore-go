package routes

import (
	"encoding/json"
	"net/http"
	sqlc "simplestore_go/internal/db"
	"simplestore_go/internal/services"
	"simplestore_go/utils"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func RegisterFilesRoutes(mux *chi.Mux, service services.Service) {
	mux.With(utils.AdminOrReadOnlyMiddleware).Route("/files", func(r chi.Router) {
		r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
			service.FileUpload(w, r)
		})
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			id := chi.URLParam(r, "id")
			var uuidId pgtype.UUID
			err := uuidId.Scan(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			file, err := service.GetFile(r.Context(), uuidId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = json.NewEncoder(w).Encode(file)
			w.WriteHeader(http.StatusCreated)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var files []sqlc.ListFilesRow
			files, err := service.GetFiles(r.Context(), sqlc.ListFilesParams{Limit: 100, Offset: 0})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = json.NewEncoder(w).Encode(files)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var file sqlc.CreateFileParams
			err := json.NewDecoder(r.Body).Decode(&file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, err = service.CreateFile(r.Context(), file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
		})
		r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			id := chi.URLParam(r, "id")
			var file sqlc.UpdateFileParams
			var uuidId pgtype.UUID
			err := uuidId.Scan(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			file.ID = uuidId
			err = json.NewDecoder(r.Body).Decode(&file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, err = service.UpdateFile(r.Context(), file)
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
			err = service.DeleteFile(r.Context(), uuidId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})
	})
}
