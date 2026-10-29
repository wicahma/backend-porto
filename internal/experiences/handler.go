package experiences

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend-porto/internal/shared/response"

	"github.com/jackc/pgx/v5/pgtype"
)

// Handler handles HTTP requests for experiences.
type Handler struct {
	svc *Service
}

// NewHandler creates a new experience Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List returns paginated experiences with nested jobs.
// GET /api/experiences?page=1&limit=20
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)
	result, err := h.svc.List(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, result)
}

// GetByID returns an experience with nested jobs.
// GET /api/experiences/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "missing id")
		return
	}

	exp, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "experience not found")
		return
	}
	response.Success(w, exp)
}

// Create creates a new experience with jobs.
// POST /api/experiences
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	exp, err := h.svc.Create(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(w, exp)
}

// Update modifies an existing experience and its jobs.
// PUT /api/experiences/{id}
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

	exp, err := h.svc.Update(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, exp)
}

// Delete removes an experience and its jobs (CASCADE).
// DELETE /api/experiences/{id}
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

// Count returns the total experience count.
// GET /api/experiences/count
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

func parseUUID(s string) (pgtype.UUID, error) {
	var uid pgtype.UUID
	if err := uid.Scan(s); err != nil {
		return uid, err
	}
	return uid, nil
}

func toPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func textToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func uuidToString(uid pgtype.UUID) string {
	if !uid.Valid {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", uid.Bytes[0:4], uid.Bytes[4:6], uid.Bytes[6:8], uid.Bytes[8:10], uid.Bytes[10:16])
}
