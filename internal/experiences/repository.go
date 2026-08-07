package experiences

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// Queries is the sqlc-generated query struct.

// Repository handles experience data persistence.
type Repository struct {
	q *Queries
}

// NewRepository creates a new Repository.
func NewRepository(q *Queries) *Repository {
	return &Repository{q: q}
}

// List returns paginated experiences.
func (r *Repository) List(ctx context.Context, limit, offset int32) ([]GetExperiencesRow, error) {
	return r.q.GetExperiences(ctx, GetExperiencesParams{Limit: limit, Offset: offset})
}

// GetByID returns a single experience by ID.
func (r *Repository) GetByID(ctx context.Context, id pgtype.UUID) (Experience, error) {
	return r.q.GetExperienceByID(ctx, id)
}

// GetJobsByExperienceID returns all jobs for an experience.
func (r *Repository) GetJobsByExperienceID(ctx context.Context, experienceID pgtype.UUID) ([]ExperienceJob, error) {
	return r.q.GetJobsByExperienceID(ctx, experienceID)
}

// GetJobsByExperienceIDs returns all jobs for multiple experiences.
func (r *Repository) GetJobsByExperienceIDs(ctx context.Context, ids []pgtype.UUID) ([]ExperienceJob, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return r.q.GetJobsByExperienceIDs(ctx, ids)
}

// CreateExperience inserts a new experience.
func (r *Repository) CreateExperience(ctx context.Context, arg CreateExperienceParams) (Experience, error) {
	return r.q.CreateExperience(ctx, arg)
}

// CreateJob inserts a new job.
func (r *Repository) CreateJob(ctx context.Context, arg CreateJobParams) (ExperienceJob, error) {
	return r.q.CreateJob(ctx, arg)
}

// UpdateExperience modifies an existing experience.
func (r *Repository) UpdateExperience(ctx context.Context, arg UpdateExperienceParams) (Experience, error) {
	return r.q.UpdateExperience(ctx, arg)
}

// DeleteJobsByExperienceID deletes all jobs for an experience.
func (r *Repository) DeleteJobsByExperienceID(ctx context.Context, experienceID pgtype.UUID) error {
	return r.q.DeleteJobsByExperienceID(ctx, experienceID)
}

// DeleteExperience removes an experience by ID.
func (r *Repository) DeleteExperience(ctx context.Context, id pgtype.UUID) error {
	return r.q.DeleteExperience(ctx, id)
}

// Count returns the total number of experiences.
func (r *Repository) Count(ctx context.Context) (int64, error) {
	return r.q.GetExperienceCount(ctx)
}
