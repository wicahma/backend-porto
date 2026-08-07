package storage

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"backend-porto/internal/shared/response"
)

// S3Storage is the interface for S3 file operations consumed by the handler.
type S3Storage interface {
	UploadFile(ctx context.Context, folder, fileName, contentType string, reader io.Reader) (string, error)
	DeleteFile(ctx context.Context, filePath string) error
	DeleteFiles(ctx context.Context, filePaths []string) error
	GetFile(ctx context.Context, filePath string) (io.ReadCloser, string, error)
}

// Handler handles HTTP requests for file storage.
type Handler struct {
	storage S3Storage
}

// NewHandler creates a new storage Handler.
func NewHandler(storage S3Storage) *Handler {
	return &Handler{storage: storage}
}

// Upload handles multipart file uploads.
// POST /api/storage/upload
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	// Limit upload size to 50MB
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "file too large or invalid form data")
		return
	}

	folder := r.FormValue("folder")
	if folder == "" {
		folder = "images"
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "missing file")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	path, err := h.storage.UploadFile(r.Context(), folder, header.Filename, contentType, file)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(w, map[string]string{"path": path})
}

// Get retrieves a file from storage.
// GET /api/storage/{*path}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	filePath := r.PathValue("path")
	if filePath == "" {
		response.Error(w, http.StatusBadRequest, "missing file path")
		return
	}

	// Remove leading slash
	filePath = strings.TrimPrefix(filePath, "/")

	body, contentType, err := h.storage.GetFile(r.Context(), filePath)
	if err != nil {
		response.Error(w, http.StatusNotFound, "file not found")
		return
	}
	defer body.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	_, _ = io.Copy(w, body)
}

// Delete removes one or more files.
// DELETE /api/storage
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Paths []string `json:"paths"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Paths) == 0 {
		response.Error(w, http.StatusBadRequest, "no file paths provided")
		return
	}

	if err := h.storage.DeleteFiles(r.Context(), req.Paths); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, map[string]string{"message": "files deleted"})
}
