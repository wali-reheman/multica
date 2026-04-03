-- name: CreateLabel :one
INSERT INTO issue_label (workspace_id, name, color)
VALUES (@workspace_id, @name, @color)
RETURNING *;

-- name: GetLabel :one
SELECT * FROM issue_label WHERE id = @id;

-- name: ListLabelsByWorkspace :many
SELECT * FROM issue_label WHERE workspace_id = @workspace_id ORDER BY name;

-- name: UpdateLabel :one
UPDATE issue_label
SET name = @name, color = @color
WHERE id = @id
RETURNING *;

-- name: DeleteLabel :exec
DELETE FROM issue_label WHERE id = @id;

-- name: AddIssueLabel :exec
INSERT INTO issue_to_label (issue_id, label_id)
VALUES (@issue_id, @label_id)
ON CONFLICT DO NOTHING;

-- name: RemoveIssueLabel :exec
DELETE FROM issue_to_label WHERE issue_id = @issue_id AND label_id = @label_id;

-- name: ListIssueLabels :many
SELECT l.* FROM issue_label l
JOIN issue_to_label il ON il.label_id = l.id
WHERE il.issue_id = @issue_id
ORDER BY l.name;

-- name: ClearIssueLabels :exec
DELETE FROM issue_to_label WHERE issue_id = @issue_id;
