-- name: GetAdminByEmail :one
SELECT id, email, github_id, created_at
FROM admin_users
WHERE email = $1;

-- name: UpsertAdmin :exec
INSERT INTO admin_users (email, github_id)
VALUES ($1, $2)
ON CONFLICT (email) DO UPDATE SET github_id = $2;
