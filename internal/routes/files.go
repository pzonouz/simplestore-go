package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	sqlc "simplestore_go/internal/db"
	"simplestore_go/internal/services"
	"simplestore_go/utils"
)

func parseUUID(w http.ResponseWriter, raw string) (pgtype.UUID, bool) {
	var id pgtype.UUID
	if err := id.Scan(raw); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return id, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// headers already sent; nothing useful left to do but log
		_ = err
	}
}

func RegisterFilesRoutes(mux *chi.Mux, service services.Service) {
	mux.With(utils.AdminOrReadOnlyMiddleware).Route("/files", func(r chi.Router) {
		r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
			service.FileUpload(w, r)
		})

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			limit, offset := int32(100), int32(0)
			if v, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && v > 0 && v <= 500 {
				limit = int32(v)
			}
			if v, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && v >= 0 {
				offset = int32(v)
			}

			files, err := service.GetFiles(
				r.Context(),
				sqlc.ListFilesParams{Limit: limit, Offset: offset},
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, http.StatusOK, files)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var file sqlc.CreateFileParams
			if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if _, err := service.CreateFile(r.Context(), file); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id, ok := parseUUID(w, chi.URLParam(r, "id"))
			if !ok {
				return
			}
			file, err := service.GetFile(r.Context(), id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					http.Error(w, "file not found", http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, http.StatusOK, file)
		})

		r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id, ok := parseUUID(w, chi.URLParam(r, "id"))
			if !ok {
				return
			}
			var file sqlc.UpdateFileParams
			if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			file.ID = id // set AFTER decoding so the body can't override it

			if _, err := service.UpdateFile(r.Context(), file); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					http.Error(w, "file not found", http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id, ok := parseUUID(w, chi.URLParam(r, "id"))
			if !ok {
				return
			}
			if err := service.DeleteFile(r.Context(), id); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})
	})
}
