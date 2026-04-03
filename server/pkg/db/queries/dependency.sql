-- name: CreateIssueDependency :one
INSERT INTO issue_dependency (issue_id, depends_on_issue_id, type)
VALUES (@issue_id, @depends_on_issue_id, @type)
RETURNING *;

-- name: DeleteIssueDependency :exec
DELETE FROM issue_dependency WHERE id = @id;

-- name: ListIssueDependencies :many
SELECT * FROM issue_dependency
WHERE issue_id = @issue_id OR depends_on_issue_id = @issue_id
ORDER BY type;
