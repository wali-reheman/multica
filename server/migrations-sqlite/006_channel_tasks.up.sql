-- Slock Phase 3: Channel-scoped agent tasks
-- Allows agents to execute tasks in response to channel messages
-- without requiring an issue. Either issue_id OR channel_id should be set.

-- SQLite cannot ALTER COLUMN to make issue_id nullable, so we recreate the table.
-- 1. Create new table with nullable issue_id + channel columns
CREATE TABLE agent_task_queue_new (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    runtime_id TEXT NOT NULL REFERENCES agent_runtime(id) ON DELETE CASCADE,
    issue_id TEXT REFERENCES issue(id) ON DELETE CASCADE,
    channel_id TEXT REFERENCES channel(id) ON DELETE CASCADE,
    channel_message_id TEXT REFERENCES channel_message(id) ON DELETE SET NULL,
    status TEXT NOT NULL DEFAULT 'queued'
        CHECK (status IN ('queued', 'dispatched', 'running', 'completed', 'failed', 'cancelled')),
    priority INTEGER NOT NULL DEFAULT 0,
    context TEXT,
    session_id TEXT,
    work_dir TEXT,
    trigger_comment_id TEXT REFERENCES comment(id) ON DELETE SET NULL,
    dispatched_at TEXT,
    started_at TEXT,
    completed_at TEXT,
    result TEXT,
    error TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

-- 2. Copy existing data
INSERT INTO agent_task_queue_new SELECT
    id, agent_id, runtime_id, issue_id,
    NULL, NULL,
    status, priority, context, session_id, work_dir,
    trigger_comment_id, dispatched_at, started_at, completed_at,
    result, error, created_at
FROM agent_task_queue;

-- 3. Drop old table and rename
DROP TABLE agent_task_queue;
ALTER TABLE agent_task_queue_new RENAME TO agent_task_queue;

-- 4. Recreate indexes
CREATE INDEX IF NOT EXISTS idx_agent_task_queue_agent ON agent_task_queue(agent_id, status);
CREATE INDEX IF NOT EXISTS idx_agent_task_queue_runtime_pending ON agent_task_queue(runtime_id, priority, created_at);
CREATE INDEX IF NOT EXISTS idx_agent_task_queue_channel ON agent_task_queue(channel_id);
