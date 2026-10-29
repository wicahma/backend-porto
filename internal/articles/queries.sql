-- name: GetArticles :many
SELECT *, COUNT(*) OVER() AS total_count
FROM articles
WHERE (@search::text = '' OR title ILIKE '%' || @search || '%' OR excerpt ILIKE '%' || @search || '%' OR category ILIKE '%' || @search || '%')
  AND (@category::text = '' OR category = @category)
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetArticleCategories :many
SELECT DISTINCT category
FROM articles
WHERE category IS NOT NULL AND category != ''
ORDER BY category ASC;

-- name: GetArticleByID :one
SELECT * FROM articles WHERE id = $1;

-- name: GetArticleBySlug :one
SELECT * FROM articles WHERE slug = $1;

-- name: GetRelatedArticles :many
SELECT * FROM articles
WHERE category = $1 AND id != $2
ORDER BY created_at DESC
LIMIT $3;

-- name: CreateArticle :one
INSERT INTO articles (
  slug, title, content, excerpt, category, image, tags, read_time,
  meta_title, meta_description, meta_keywords, og_image
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: UpdateArticle :one
UPDATE articles SET
  slug = $2,
  title = $3,
  content = $4,
  excerpt = $5,
  category = $6,
  image = $7,
  tags = $8,
  read_time = $9,
  meta_title = $10,
  meta_description = $11,
  meta_keywords = $12,
  og_image = $13,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteArticle :exec
DELETE FROM articles WHERE id = $1;

-- name: CheckSlugExists :one
SELECT EXISTS(
  SELECT 1 FROM articles WHERE slug = $1 AND ($2::uuid IS NULL OR id != $2)
) AS exists;

-- name: GetArticleCount :one
SELECT COUNT(*) FROM articles;
