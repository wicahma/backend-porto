package articles

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"backend-porto/internal/shared/response"

	"github.com/jackc/pgx/v5/pgtype"
)

// Handler handles HTTP requests for articles.
type Handler struct {
	svc *Service
}

// NewHandler creates a new article Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List returns paginated articles.
// GET /api/articles?page=1&limit=20
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)
	result, err := h.svc.List(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, result)
}

// GetByID returns an article by ID.
// GET /api/articles/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "missing id")
		return
	}

	article, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "article not found")
		return
	}
	response.Success(w, article)
}

// GetBySlug returns an article by slug.
// GET /api/articles/slug/{slug}
func (h *Handler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		response.Error(w, http.StatusBadRequest, "missing slug")
		return
	}

	article, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		response.Error(w, http.StatusNotFound, "article not found")
		return
	}
	response.Success(w, article)
}

// GetRelated returns related articles.
// GET /api/articles/{id}/related?limit=5
func (h *Handler) GetRelated(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "missing id")
		return
	}

	article, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "article not found")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 5
	}

	related, err := h.svc.GetRelated(r.Context(), id, article.Category, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, related)
}

// Create creates a new article.
// POST /api/articles
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	article, err := h.svc.Create(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(w, article)
}

// Update modifies an existing article.
// PUT /api/articles/{id}
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

	article, err := h.svc.Update(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, article)
}

// Delete removes an article.
// DELETE /api/articles/{id}
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

// Count returns the total article count.
// GET /api/articles/count
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
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
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

func randomSuffix() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
