-- name: CreateChannel :one
INSERT INTO channel (id, workspace_id, name, description, type, created_by_type, created_by_id)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetChannel :one
SELECT * FROM channel WHERE id = ?;

-- name: GetChannelInWorkspace :one
SELECT * FROM channel WHERE id = ? AND workspace_id = ?;

-- name: ListChannels :many
SELECT * FROM channel
WHERE workspace_id = ? AND archived_at IS NULL
ORDER BY name ASC;

-- name: UpdateChannel :one
UPDATE channel SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: ArchiveChannel :exec
UPDATE channel SET archived_at = datetime('now'), updated_at = datetime('now') WHERE id = ?;

-- name: DeleteChannel :exec
DELETE FROM channel WHERE id = ?;

-- name: AddChannelMember :exec
INSERT INTO channel_member (channel_id, member_type, member_id, role)
VALUES (?, ?, ?, ?)
ON CONFLICT (channel_id, member_type, member_id) DO NOTHING;

-- name: RemoveChannelMember :exec
DELETE FROM channel_member
WHERE channel_id = ? AND member_type = ? AND member_id = ?;

-- name: ListChannelMembers :many
SELECT * FROM channel_member WHERE channel_id = ?;

-- name: ListChannelsForMember :many
SELECT c.* FROM channel c
JOIN channel_member cm ON cm.channel_id = c.id
WHERE c.workspace_id = ? AND cm.member_type = ? AND cm.member_id = ? AND c.archived_at IS NULL
ORDER BY c.updated_at DESC;

-- name: IsChannelMember :one
SELECT count(*) FROM channel_member
WHERE channel_id = ? AND member_type = ? AND member_id = ?;

-- name: CreateChannelMessage :one
INSERT INTO channel_message (id, channel_id, workspace_id, author_type, author_id, content, type, parent_id, issue_id)
VALUES (?, ?, ?, ?, ?, ?, ?, sqlc.narg(parent_id), sqlc.narg(issue_id))
RETURNING *;

-- name: GetChannelMessage :one
SELECT * FROM channel_message WHERE id = ?;

-- name: ListChannelMessages :many
SELECT * FROM channel_message
WHERE channel_id = ? AND workspace_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: ListChannelMessagesSince :many
SELECT * FROM channel_message
WHERE channel_id = ? AND workspace_id = ? AND created_at > ?
ORDER BY created_at ASC
LIMIT ?;

-- name: CountChannelMessages :one
SELECT count(*) FROM channel_message
WHERE channel_id = ? AND workspace_id = ?;

-- name: UpdateChannelMessage :one
UPDATE channel_message SET
    content = ?,
    updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: DeleteChannelMessage :exec
DELETE FROM channel_message WHERE id = ?;

-- name: TouchChannel :exec
UPDATE channel SET updated_at = datetime('now') WHERE id = ?;
