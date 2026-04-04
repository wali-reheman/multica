-- Slock: Slack-like channels for group chat between agents and members
CREATE TABLE IF NOT EXISTS channel (
    id TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL REFERENCES workspace(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    type TEXT NOT NULL DEFAULT 'group' CHECK (type IN ('group', 'direct')),
    created_by_type TEXT NOT NULL CHECK (created_by_type IN ('member', 'agent')),
    created_by_id TEXT NOT NULL,
    archived_at TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    UNIQUE(workspace_id, name)
);

CREATE TABLE IF NOT EXISTS channel_member (
    channel_id TEXT NOT NULL REFERENCES channel(id) ON DELETE CASCADE,
    member_type TEXT NOT NULL CHECK (member_type IN ('member', 'agent')),
    member_id TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'member')),
    joined_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    PRIMARY KEY (channel_id, member_type, member_id)
);

CREATE TABLE IF NOT EXISTS channel_message (
    id TEXT PRIMARY KEY,
    channel_id TEXT NOT NULL REFERENCES channel(id) ON DELETE CASCADE,
    workspace_id TEXT NOT NULL REFERENCES workspace(id) ON DELETE CASCADE,
    author_type TEXT NOT NULL CHECK (author_type IN ('member', 'agent')),
    author_id TEXT NOT NULL,
    content TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT 'message' CHECK (type IN ('message', 'system', 'issue_created')),
    parent_id TEXT REFERENCES channel_message(id) ON DELETE CASCADE,
    issue_id TEXT REFERENCES issue(id) ON DELETE SET NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_channel_workspace ON channel(workspace_id);
CREATE INDEX IF NOT EXISTS idx_channel_member_channel ON channel_member(channel_id);
CREATE INDEX IF NOT EXISTS idx_channel_member_member ON channel_member(member_type, member_id);
CREATE INDEX IF NOT EXISTS idx_channel_message_channel ON channel_message(channel_id, created_at);
CREATE INDEX IF NOT EXISTS idx_channel_message_workspace ON channel_message(workspace_id);
CREATE INDEX IF NOT EXISTS idx_channel_message_parent ON channel_message(parent_id);
