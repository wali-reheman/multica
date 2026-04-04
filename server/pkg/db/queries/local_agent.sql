-- MULTICA-LOCAL: Stage 4 - Local Agent Config & Skills queries (SQLite dialect)

-- name: ListLocalAgentConfigs :many
SELECT * FROM local_agent_config
WHERE workspace_id = ?
ORDER BY provider ASC;

-- name: GetLocalAgentConfig :one
SELECT * FROM local_agent_config
WHERE workspace_id = ? AND provider = ?;

-- name: UpsertLocalAgentConfig :one
INSERT INTO local_agent_config (id, workspace_id, provider, cli_path, version, status, is_custom_path, last_health_check, health_error)
VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), ?)
ON CONFLICT (workspace_id, provider)
DO UPDATE SET
    cli_path = excluded.cli_path,
    version = excluded.version,
    status = excluded.status,
    is_custom_path = excluded.is_custom_path,
    last_health_check = datetime('now'),
    health_error = excluded.health_error,
    updated_at = datetime('now')
RETURNING *;

-- name: ListLocalSkills :many
SELECT * FROM local_skill
WHERE (workspace_id = ? OR workspace_id IS NULL)
  AND (project_path = ? OR project_path IS NULL)
ORDER BY is_default DESC, name ASC;

-- name: ListGlobalLocalSkills :many
SELECT * FROM local_skill
WHERE workspace_id IS NULL
ORDER BY name ASC;

-- name: GetLocalSkill :one
SELECT * FROM local_skill WHERE id = ?;

-- name: CreateLocalSkill :one
INSERT INTO local_skill (id, workspace_id, project_path, name, description, content, is_default)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateLocalSkill :one
UPDATE local_skill SET
    name = COALESCE(NULLIF(?, ''), name),
    description = COALESCE(NULLIF(?, ''), description),
    content = COALESCE(NULLIF(?, ''), content),
    updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: DeleteLocalSkill :exec
DELETE FROM local_skill WHERE id = ?;
