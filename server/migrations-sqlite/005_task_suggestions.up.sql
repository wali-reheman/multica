-- Slock Phase 2: Task suggestions with approval flow
-- Allows members and agents to propose tasks in channels,
-- which can be approved (creating an issue) or dismissed.
CREATE TABLE IF NOT EXISTS task_suggestion (
    id TEXT PRIMARY KEY,
    channel_id TEXT NOT NULL REFERENCES channel(id) ON DELETE CASCADE,
    workspace_id TEXT NOT NULL REFERENCES workspace(id) ON DELETE CASCADE,
    message_id TEXT REFERENCES channel_message(id) ON DELETE SET NULL,
    suggested_by_type TEXT NOT NULL CHECK (suggested_by_type IN ('member', 'agent')),
    suggested_by_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    priority TEXT NOT NULL DEFAULT 'none'
        CHECK (priority IN ('urgent', 'high', 'medium', 'low', 'none')),
    assignee_type TEXT CHECK (assignee_type IN ('member', 'agent')),
    assignee_id TEXT,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'approved', 'dismissed')),
    resolved_by_type TEXT CHECK (resolved_by_type IN ('member', 'agent')),
    resolved_by_id TEXT,
    issue_id TEXT REFERENCES issue(id) ON DELETE SET NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_task_suggestion_channel ON task_suggestion(channel_id, status);
CREATE INDEX IF NOT EXISTS idx_task_suggestion_workspace ON task_suggestion(workspace_id);
CREATE INDEX IF NOT EXISTS idx_task_suggestion_issue ON task_suggestion(issue_id);
