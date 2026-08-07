-- name: GetExperiences :many
SELECT *, COUNT(*) OVER() AS total_count
FROM experiences
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetExperienceByID :one
SELECT * FROM experiences WHERE id = $1;

-- name: GetJobsByExperienceID :many
SELECT * FROM experience_jobs
WHERE experience_id = $1
ORDER BY start_date DESC;

-- name: GetJobsByExperienceIDs :many
SELECT * FROM experience_jobs
WHERE experience_id = ANY($1::uuid[])
ORDER BY start_date DESC;

-- name: CreateExperience :one
INSERT INTO experiences (company, location, total_duration, tags, is_current)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateJob :one
INSERT INTO experience_jobs (experience_id, position, employment_type, description, competency, start_date, end_date, is_current)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateExperience :one
UPDATE experiences SET
  company = $2,
  location = $3,
  total_duration = $4,
  tags = $5,
  is_current = $6,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteJobsByExperienceID :exec
DELETE FROM experience_jobs WHERE experience_id = $1;

-- name: DeleteExperience :exec
DELETE FROM experiences WHERE id = $1;

-- name: GetExperienceCount :one
SELECT COUNT(*) FROM experiences;
