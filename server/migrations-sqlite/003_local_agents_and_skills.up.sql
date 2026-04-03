-- Stage 4: Local agent configuration and local skills
CREATE TABLE IF NOT EXISTS local_agent_config (
    id TEXT PRIMARY KEY,
    workspace_id TEXT NOT NULL REFERENCES workspace(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    cli_path TEXT NOT NULL DEFAULT '',
    version TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'unknown',
    is_custom_path INTEGER NOT NULL DEFAULT 0,
    last_health_check TEXT,
    health_error TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    UNIQUE(workspace_id, provider)
);

CREATE TABLE IF NOT EXISTS local_skill (
    id TEXT PRIMARY KEY,
    workspace_id TEXT REFERENCES workspace(id) ON DELETE CASCADE,
    project_path TEXT,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    content TEXT NOT NULL DEFAULT '',
    is_default INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);
