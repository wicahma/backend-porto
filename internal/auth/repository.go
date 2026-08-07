package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// Queries is the sqlc-generated query struct. We alias it for clarity.

// Repository handles admin user data persistence.
type Repository struct {
	q *Queries
}

// NewRepository creates a new Repository.
func NewRepository(q *Queries) *Repository {
	return &Repository{q: q}
}

// GetAdminByEmail retrieves an admin user by email.
func (r *Repository) GetAdminByEmail(ctx context.Context, email string) (AdminUser, error) {
	return r.q.GetAdminByEmail(ctx, email)
}

// UpsertAdmin creates or updates an admin user's GitHub ID by email.
func (r *Repository) UpsertAdmin(ctx context.Context, email, githubID string) error {
	return r.q.UpsertAdmin(ctx, UpsertAdminParams{
		Email:    email,
		GithubID: pgtype.Text{String: githubID, Valid: githubID != ""},
	})
}
