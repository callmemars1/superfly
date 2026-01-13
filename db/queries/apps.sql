-- name: CreateApp :one
INSERT INTO apps (
    slug,
    name,
    image,
    port,
    replicas,
    cpu_limit,
    memory_limit,
    domain,
    health_check_path,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetApp :one
SELECT * FROM apps
WHERE id = $1 LIMIT 1;

-- name: GetAppBySlug :one
SELECT * FROM apps
WHERE slug = $1 LIMIT 1;

-- name: ListApps :many
SELECT * FROM apps
ORDER BY created_at DESC;

-- name: UpdateAppStatus :one
UPDATE apps
SET status = $2,
    updated_at = NOW(),
    last_deployed_at = CASE WHEN $2 = 'running' THEN NOW() ELSE last_deployed_at END
WHERE id = $1
RETURNING *;

-- name: UpdateApp :one
UPDATE apps
SET name = COALESCE($2, name),
    image = COALESCE($3, image),
    port = COALESCE($4, port),
    replicas = COALESCE($5, replicas),
    cpu_limit = COALESCE($6, cpu_limit),
    memory_limit = COALESCE($7, memory_limit),
    domain = COALESCE($8, domain),
    health_check_path = COALESCE($9, health_check_path),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteApp :exec
DELETE FROM apps
WHERE id = $1;

-- name: CheckSlugExists :one
SELECT EXISTS(SELECT 1 FROM apps WHERE slug = $1);

-- name: CheckDomainExists :one
SELECT EXISTS(SELECT 1 FROM apps WHERE domain = $1);
