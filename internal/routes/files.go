package routes

import (
	"net/http"
	"simplestore_go/internal/services"

	"github.com/go-chi/chi/v5"
)

func RegisterFilesRoutes(mux *chi.Mux, service services.Service) {
	mux.Route("/files", func(r chi.Router) {
		r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
			service.FileUpload(w, r)
		})
	})
}
