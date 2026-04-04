-- name: CreateTaskSuggestion :one
INSERT INTO task_suggestion (
    id, channel_id, workspace_id, message_id,
    suggested_by_type, suggested_by_id,
    title, description, priority,
    assignee_type, assignee_id
) VALUES (
    ?, ?, ?, sqlc.narg(message_id),
    ?, ?,
    ?, ?, ?,
    sqlc.narg(assignee_type), sqlc.narg(assignee_id)
) RETURNING *;

-- name: GetTaskSuggestion :one
SELECT * FROM task_suggestion WHERE id = ?;

-- name: GetTaskSuggestionInWorkspace :one
SELECT * FROM task_suggestion WHERE id = ? AND workspace_id = ?;

-- name: ListTaskSuggestions :many
SELECT * FROM task_suggestion
WHERE channel_id = ? AND workspace_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListPendingTaskSuggestions :many
SELECT * FROM task_suggestion
WHERE channel_id = ? AND workspace_id = ? AND status = 'pending'
ORDER BY created_at DESC;

-- name: ApproveTaskSuggestion :one
UPDATE task_suggestion SET
    status = 'approved',
    resolved_by_type = ?,
    resolved_by_id = ?,
    issue_id = ?,
    updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: DismissTaskSuggestion :one
UPDATE task_suggestion SET
    status = 'dismissed',
    resolved_by_type = ?,
    resolved_by_id = ?,
    updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: CountPendingSuggestions :one
SELECT count(*) FROM task_suggestion
WHERE channel_id = ? AND workspace_id = ? AND status = 'pending';

-- name: UpdateTaskSuggestionMessage :exec
UPDATE task_suggestion SET message_id = sqlc.narg(message_id) WHERE id = ?;
