-- name: CreateLabel :one
INSERT INTO issue_label (id, workspace_id, name, color)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetLabel :one
SELECT * FROM issue_label WHERE id = ?;

-- name: ListLabelsByWorkspace :many
SELECT * FROM issue_label WHERE workspace_id = ? ORDER BY name;

-- name: UpdateLabel :one
UPDATE issue_label
SET name = ?, color = ?
WHERE id = ?
RETURNING *;

-- name: DeleteLabel :exec
DELETE FROM issue_label WHERE id = ?;

-- name: AddIssueLabel :exec
INSERT OR IGNORE INTO issue_to_label (issue_id, label_id)
VALUES (?, ?);

-- name: RemoveIssueLabel :exec
DELETE FROM issue_to_label WHERE issue_id = ? AND label_id = ?;

-- name: ListIssueLabels :many
SELECT l.* FROM issue_label l
JOIN issue_to_label il ON il.label_id = l.id
WHERE il.issue_id = ?
ORDER BY l.name;

-- name: ClearIssueLabels :exec
DELETE FROM issue_to_label WHERE issue_id = ?;
