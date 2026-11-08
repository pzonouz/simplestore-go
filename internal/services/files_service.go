package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	sqlc "simplestore_go/internal/db"
	"strings"
	"uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) FileUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	ext := filepath.Ext(header.Filename)
	safeName := (strings.ReplaceAll(uuid.New().String(), "-", "")[:16] + ext)
	dst, err := os.Create(filepath.Join("uploads", safeName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	type UploadReturnData struct {
		Url string `json:"url"`
	}
	data := UploadReturnData{
		Url: dst.Name(),
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (s *Service) GetFiles(ctx context.Context, args sqlc.ListFilesParams) ([]sqlc.ListFilesRow, error) {
	Files, err := s.queries.ListFiles(ctx, args)
	if err != nil {
		return []sqlc.ListFilesRow{}, err
	}
	if Files == nil {
		return []sqlc.ListFilesRow{}, err
	}
	return Files, err
}
func (s *Service) GetFile(ctx context.Context, id pgtype.UUID) (sqlc.GetFileRow, error) {

	File, err := s.queries.GetFile(ctx, id)
	if err != nil {
		return sqlc.GetFileRow{}, err
	}
	return File, err
}
func (s *Service) CreateFile(ctx context.Context, args sqlc.CreateFileParams) (sqlc.File, error) {
	file, err := s.queries.CreateFile(ctx, args)
	if err != nil {
		return sqlc.File{}, err
	}

	return file, err
}
func (s *Service) UpdateFile(ctx context.Context, args sqlc.UpdateFileParams) (sqlc.File, error) {
	return s.queries.UpdateFile(ctx, args)
}
func (s *Service) DeleteFile(ctx context.Context, id pgtype.UUID) error {
	return s.queries.DeleteFile(ctx, id)
}
