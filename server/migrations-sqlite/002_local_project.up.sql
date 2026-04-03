-- Local projects: user-tracked local folders with git version history.
CREATE TABLE IF NOT EXISTS local_project (
    id TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL REFERENCES workspace(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    local_path TEXT NOT NULL,
    default_branch TEXT NOT NULL DEFAULT 'main',
    language TEXT,
    file_count INTEGER NOT NULL DEFAULT 0,
    size_bytes INTEGER NOT NULL DEFAULT 0,
    last_opened_at TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    UNIQUE(workspace_id, local_path)
);

CREATE INDEX IF NOT EXISTS idx_local_project_workspace ON local_project(workspace_id);
