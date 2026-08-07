package articles

import (
	"context"
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5/pgtype"
)

// CreateInput is the input for creating an article.
type CreateInput struct {
	Title           string
	Content         string
	Excerpt         string
	Category        string
	Image           string
	Tags            []string
	ReadTime        string
	MetaTitle       string
	MetaDescription string
	MetaKeywords    []string
	OgImage         string
}

// UpdateInput is the input for updating an article.
type UpdateInput struct {
	ID              string
	Title           string
	Content         string
	Excerpt         string
	Category        string
	Image           string
	Tags            []string
	ReadTime        string
	MetaTitle       string
	MetaDescription string
	MetaKeywords    []string
	OgImage         string
}

// PaginatedResponse represents a paginated list of articles.
type PaginatedResponse struct {
	Articles   []Article `json:"articles"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	TotalPages int       `json:"total_pages"`
}

// Service handles article business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new article Service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List returns paginated articles.
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
		return nil, fmt.Errorf("listing articles: %w", err)
	}

	articles := make([]Article, len(rows))
	total := int64(0)
	for i, row := range rows {
		articles[i] = Article{
			ID:              row.ID,
			Slug:            row.Slug,
			Title:           row.Title,
			Content:         row.Content,
			Excerpt:         row.Excerpt,
			Category:        row.Category,
			Image:           row.Image,
			Tags:            row.Tags,
			ReadTime:        row.ReadTime,
			MetaTitle:       row.MetaTitle,
			MetaDescription: row.MetaDescription,
			MetaKeywords:    row.MetaKeywords,
			OgImage:         row.OgImage,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
		}
		total = row.TotalCount
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &PaginatedResponse{
		Articles:   articles,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	}, nil
}

// GetByID returns an article by ID.
func (s *Service) GetByID(ctx context.Context, id string) (Article, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return Article{}, fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.GetByID(ctx, uid)
}

// GetBySlug returns an article by slug.
func (s *Service) GetBySlug(ctx context.Context, slug string) (Article, error) {
	return s.repo.GetBySlug(ctx, slug)
}

// GetRelated returns related articles by category.
func (s *Service) GetRelated(ctx context.Context, id, category string, limit int) ([]Article, error) {
	if limit < 1 || limit > 20 {
		limit = 5
	}
	uid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.GetRelated(ctx, category, uid, int32(limit))
}

// Create creates a new article with slug generation and SEO defaults.
func (s *Service) Create(ctx context.Context, input CreateInput) (Article, error) {
	slug := generateSlug(input.Title)

	// Validate slug uniqueness
	exists, err := s.repo.CheckSlugExists(ctx, slug, pgtype.UUID{Valid: false})
	if err != nil {
		return Article{}, fmt.Errorf("checking slug: %w", err)
	}
	if exists {
		slug = slug + "-" + randomSuffix()
	}

	// SEO defaults
	metaTitle := toPgText(input.MetaTitle)
	if !metaTitle.Valid || metaTitle.String == "" {
		metaTitle = pgtype.Text{String: input.Title, Valid: true}
	}

	metaDescription := toPgText(input.MetaDescription)
	if !metaDescription.Valid || metaDescription.String == "" {
		metaDescription = toPgText(input.Excerpt)
	}

	ogImage := toPgText(input.OgImage)
	if !ogImage.Valid || ogImage.String == "" {
		ogImage = pgtype.Text{String: input.Image, Valid: true}
	}

	return s.repo.Create(ctx, CreateArticleParams{
		Slug:            slug,
		Title:           input.Title,
		Content:         input.Content,
		Excerpt:         toPgText(input.Excerpt),
		Category:        input.Category,
		Image:           input.Image,
		Tags:            input.Tags,
		ReadTime:        input.ReadTime,
		MetaTitle:       metaTitle,
		MetaDescription: metaDescription,
		MetaKeywords:    input.MetaKeywords,
		OgImage:         ogImage,
	})
}

// Update modifies an existing article.
func (s *Service) Update(ctx context.Context, input UpdateInput) (Article, error) {
	uid, err := parseUUID(input.ID)
	if err != nil {
		return Article{}, fmt.Errorf("invalid id: %w", err)
	}

	slug := generateSlug(input.Title)

	// Check slug uniqueness excluding self
	exists, err := s.repo.CheckSlugExists(ctx, slug, uid)
	if err != nil {
		return Article{}, fmt.Errorf("checking slug: %w", err)
	}
	if exists {
		slug = slug + "-" + randomSuffix()
	}

	// SEO defaults
	metaTitle := toPgText(input.MetaTitle)
	if !metaTitle.Valid || metaTitle.String == "" {
		metaTitle = pgtype.Text{String: input.Title, Valid: true}
	}

	metaDescription := toPgText(input.MetaDescription)
	if !metaDescription.Valid || metaDescription.String == "" {
		metaDescription = toPgText(input.Excerpt)
	}

	ogImage := toPgText(input.OgImage)
	if !ogImage.Valid || ogImage.String == "" {
		ogImage = pgtype.Text{String: input.Image, Valid: true}
	}

	return s.repo.Update(ctx, UpdateArticleParams{
		ID:              uid,
		Slug:            slug,
		Title:           input.Title,
		Content:         input.Content,
		Excerpt:         toPgText(input.Excerpt),
		Category:        input.Category,
		Image:           input.Image,
		Tags:            input.Tags,
		ReadTime:        input.ReadTime,
		MetaTitle:       metaTitle,
		MetaDescription: metaDescription,
		MetaKeywords:    input.MetaKeywords,
		OgImage:         ogImage,
	})
}

// Delete removes an article by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.Delete(ctx, uid)
}

// Count returns the total number of articles.
func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.TrimSpace(slug)

	var result strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' {
			if r == ' ' {
				result.WriteRune('-')
			} else {
				result.WriteRune(r)
			}
		}
	}

	s := result.String()
	// Remove consecutive dashes
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")

	if s == "" {
		s = "untitled"
	}

	return s
}
