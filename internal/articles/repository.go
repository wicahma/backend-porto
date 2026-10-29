package articles

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// Queries is the sqlc-generated query struct.

// Repository handles article data persistence.
type Repository struct {
	q *Queries
}

// NewRepository creates a new Repository.
func NewRepository(q *Queries) *Repository {
	return &Repository{q: q}
}

// List returns paginated articles ordered by created_at DESC with optional search and category filter.
func (r *Repository) List(ctx context.Context, limit, offset int32, search, category string) ([]GetArticlesRow, error) {
	return r.q.GetArticles(ctx, GetArticlesParams{
		Limit:    limit,
		Offset:   offset,
		Search:   search,
		Category: category,
	})
}

// GetCategories returns all distinct article categories.
func (r *Repository) GetCategories(ctx context.Context) ([]string, error) {
	return r.q.GetArticleCategories(ctx)
}
// GetByID returns a single article by ID.
func (r *Repository) GetByID(ctx context.Context, id pgtype.UUID) (Article, error) {
	return r.q.GetArticleByID(ctx, id)
}

// GetBySlug returns a single article by slug.
func (r *Repository) GetBySlug(ctx context.Context, slug string) (Article, error) {
	return r.q.GetArticleBySlug(ctx, slug)
}

// GetRelated returns related articles by category excluding the given ID.
func (r *Repository) GetRelated(ctx context.Context, category string, excludeID pgtype.UUID, limit int32) ([]Article, error) {
	return r.q.GetRelatedArticles(ctx, GetRelatedArticlesParams{
		Category: category,
		ID:       excludeID,
		Limit:    limit,
	})
}

// Create inserts a new article.
func (r *Repository) Create(ctx context.Context, arg CreateArticleParams) (Article, error) {
	return r.q.CreateArticle(ctx, arg)
}

// Update modifies an existing article.
func (r *Repository) Update(ctx context.Context, arg UpdateArticleParams) (Article, error) {
	return r.q.UpdateArticle(ctx, arg)
}

// Delete removes an article by ID.
func (r *Repository) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.q.DeleteArticle(ctx, id)
}

// CheckSlugExists checks if a slug exists, optionally excluding a specific ID.
func (r *Repository) CheckSlugExists(ctx context.Context, slug string, excludeID pgtype.UUID) (bool, error) {
	return r.q.CheckSlugExists(ctx, CheckSlugExistsParams{
		Slug:    slug,
		Column2: excludeID,
	})
}

// Count returns the total number of articles.
func (r *Repository) Count(ctx context.Context) (int64, error) {
	return r.q.GetArticleCount(ctx)
}
