package projects

import (
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

// CreateInput is the input for creating a project.
type CreateInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Problem     string   `json:"problem,omitempty"`
	Solution    string   `json:"solution,omitempty"`
	Image       string   `json:"image"`
	Tags        []string `json:"tags"`
	Year        string   `json:"year"`
	Link        string   `json:"link,omitempty"`
}

// UpdateInput is the input for updating a project.
type UpdateInput struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Problem     string   `json:"problem,omitempty"`
	Solution    string   `json:"solution,omitempty"`
	Image       string   `json:"image"`
	Tags        []string `json:"tags"`
	Year        string   `json:"year"`
	Link        string   `json:"link,omitempty"`
}

// PaginatedResponse represents a paginated list of projects.
type PaginatedResponse struct {
	Projects   []Project `json:"projects"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	TotalPages int       `json:"total_pages"`
}

// Service handles project business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new project Service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List returns paginated projects.
func (s *Service) List(ctx context.Context, page, limit int) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	rows, err := s.repo.List(ctx, int32(limit), int32(offset))
	if err != nil {
		return nil, fmt.Errorf("listing projects: %w", err)
	}

	projects := make([]Project, len(rows))
	total := int64(0)
	for i, row := range rows {
		projects[i] = Project{
			ID:          row.ID,
			Title:       row.Title,
			Description: row.Description,
			Problem:     row.Problem,
			Solution:    row.Solution,
			Image:       row.Image,
			Tags:        row.Tags,
			Year:        row.Year,
			Link:        row.Link,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		}
		total = row.TotalCount
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &PaginatedResponse{
		Projects:   projects,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	}, nil
}

// GetByID returns a project by ID.
func (s *Service) GetByID(ctx context.Context, id string) (Project, error) {
	var uid pgtype.UUID
	if err := uid.Scan(id); err != nil {
		return Project{}, fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.GetByID(ctx, uid)
}

// Create creates a new project.
func (s *Service) Create(ctx context.Context, input CreateInput) (Project, error) {
	return s.repo.Create(ctx, CreateProjectParams{
		Title:       input.Title,
		Description: input.Description,
		Problem:     toPgText(input.Problem),
		Solution:    toPgText(input.Solution),
		Image:       input.Image,
		Tags:        input.Tags,
		Year:        input.Year,
		Link:        toPgText(input.Link),
	})
}

// Update modifies an existing project.
func (s *Service) Update(ctx context.Context, input UpdateInput) (Project, error) {
	var uid pgtype.UUID
	if err := uid.Scan(input.ID); err != nil {
		return Project{}, fmt.Errorf("invalid id: %w", err)
	}

	return s.repo.Update(ctx, UpdateProjectParams{
		ID:          uid,
		Title:       input.Title,
		Description: input.Description,
		Problem:     toPgText(input.Problem),
		Solution:    toPgText(input.Solution),
		Image:       input.Image,
		Tags:        input.Tags,
		Year:        input.Year,
		Link:        toPgText(input.Link),
	})
}

// Delete removes a project by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	var uid pgtype.UUID
	if err := uid.Scan(id); err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.Delete(ctx, uid)
}

// Count returns the total number of projects.
func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func toPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}
