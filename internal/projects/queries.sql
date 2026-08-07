-- name: GetProjects :many
SELECT *, COUNT(*) OVER() AS total_count
FROM projects
ORDER BY year DESC
LIMIT $1 OFFSET $2;

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id = $1;

-- name: CreateProject :one
INSERT INTO projects (title, description, problem, solution, image, tags, year, link)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateProject :one
UPDATE projects SET
  title = $2,
  description = $3,
  problem = $4,
  solution = $5,
  image = $6,
  tags = $7,
  year = $8,
  link = $9,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;

-- name: GetProjectCount :one
SELECT COUNT(*) FROM projects;
