package projects

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend-porto/internal/shared/response"
)

// Handler handles HTTP requests for projects.
type Handler struct {
	svc *Service
}

// NewHandler creates a new project Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List returns paginated projects.
// GET /api/projects?page=1&limit=20
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)
	result, err := h.svc.List(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, result)
}

// GetByID returns a project by ID.
// GET /api/projects/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "missing id")
		return
	}

	project, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "project not found")
		return
	}
	response.Success(w, project)
}

// Create creates a new project.
// POST /api/projects
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := h.svc.Create(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(w, project)
}

// Update modifies an existing project.
// PUT /api/projects/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "missing id")
		return
	}

	var input UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input.ID = id

	project, err := h.svc.Update(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, project)
}

// Delete removes a project.
// DELETE /api/projects/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "missing id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.NoContent(w)
}

// Count returns the total project count.
// GET /api/projects/count
func (h *Handler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.svc.Count(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, map[string]int64{"count": count})
}

// --- helpers ---

func parsePagination(r *http.Request) (int, int) {
	page := 1
	limit := 20
	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}
