-- name: ListLocalProjects :many
SELECT * FROM local_project
WHERE workspace_id = ?
ORDER BY last_opened_at DESC, created_at DESC
LIMIT ? OFFSET ?;

-- name: GetLocalProject :one
SELECT * FROM local_project
WHERE id = ?;

-- name: GetLocalProjectInWorkspace :one
SELECT * FROM local_project
WHERE id = ? AND workspace_id = ?;

-- name: GetLocalProjectByPath :one
SELECT * FROM local_project
WHERE workspace_id = ? AND local_path = ?;

-- name: CreateLocalProject :one
INSERT INTO local_project (
    id, workspace_id, name, local_path, default_branch,
    language, file_count, size_bytes
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateLocalProject :one
UPDATE local_project SET
    name = COALESCE(sqlc.narg('name'), name),
    default_branch = COALESCE(sqlc.narg('default_branch'), default_branch),
    language = COALESCE(sqlc.narg('language'), language),
    file_count = COALESCE(sqlc.narg('file_count'), file_count),
    size_bytes = COALESCE(sqlc.narg('size_bytes'), size_bytes),
    updated_at = strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
WHERE id = ?
RETURNING *;

-- name: UpdateLocalProjectLastOpened :exec
UPDATE local_project SET
    last_opened_at = strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
WHERE id = ?;

-- name: DeleteLocalProject :exec
DELETE FROM local_project WHERE id = ?;
