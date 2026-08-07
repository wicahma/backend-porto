package projects

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// Queries is the sqlc-generated query struct.

// Repository handles project data persistence.
type Repository struct {
	q *Queries
}

// NewRepository creates a new Repository.
func NewRepository(q *Queries) *Repository {
	return &Repository{q: q}
}

// List returns paginated projects ordered by year DESC.
func (r *Repository) List(ctx context.Context, limit, offset int32) ([]GetProjectsRow, error) {
	return r.q.GetProjects(ctx, GetProjectsParams{Limit: limit, Offset: offset})
}

// GetByID returns a single project by ID.
func (r *Repository) GetByID(ctx context.Context, id pgtype.UUID) (Project, error) {
	return r.q.GetProjectByID(ctx, id)
}

// Create inserts a new project.
func (r *Repository) Create(ctx context.Context, arg CreateProjectParams) (Project, error) {
	return r.q.CreateProject(ctx, arg)
}

// Update modifies an existing project.
func (r *Repository) Update(ctx context.Context, arg UpdateProjectParams) (Project, error) {
	return r.q.UpdateProject(ctx, arg)
}

// Delete removes a project by ID.
func (r *Repository) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.q.DeleteProject(ctx, id)
}

// Count returns the total number of projects.
func (r *Repository) Count(ctx context.Context) (int64, error) {
	return r.q.GetProjectCount(ctx)
}
