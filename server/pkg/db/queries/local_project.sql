-- name: ListLocalProjects :many
SELECT * FROM local_project
WHERE workspace_id = $1
ORDER BY last_opened_at DESC NULLS LAST, created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetLocalProject :one
SELECT * FROM local_project
WHERE id = $1;

-- name: GetLocalProjectInWorkspace :one
SELECT * FROM local_project
WHERE id = $1 AND workspace_id = $2;

-- name: GetLocalProjectByPath :one
SELECT * FROM local_project
WHERE workspace_id = $1 AND local_path = $2;

-- name: CreateLocalProject :one
INSERT INTO local_project (
    workspace_id, name, local_path, default_branch,
    language, file_count, size_bytes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdateLocalProject :one
UPDATE local_project SET
    name = COALESCE(sqlc.narg('name'), name),
    default_branch = COALESCE(sqlc.narg('default_branch'), default_branch),
    language = COALESCE(sqlc.narg('language'), language),
    file_count = COALESCE(sqlc.narg('file_count'), file_count),
    size_bytes = COALESCE(sqlc.narg('size_bytes'), size_bytes),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateLocalProjectLastOpened :exec
UPDATE local_project SET
    last_opened_at = now()
WHERE id = $1;

-- name: DeleteLocalProject :exec
DELETE FROM local_project WHERE id = $1;
