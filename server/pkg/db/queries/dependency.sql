-- name: CreateIssueDependency :one
INSERT INTO issue_dependency (id, issue_id, depends_on_issue_id, type)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteIssueDependency :exec
DELETE FROM issue_dependency WHERE id = ?;

-- name: ListIssueDependencies :many
SELECT * FROM issue_dependency
WHERE issue_id = ? OR depends_on_issue_id = ?
ORDER BY type;
