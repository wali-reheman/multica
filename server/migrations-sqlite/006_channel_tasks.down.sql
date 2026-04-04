-- Revert: restore agent_task_queue with NOT NULL issue_id
CREATE TABLE agent_task_queue_old (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    runtime_id TEXT NOT NULL REFERENCES agent_runtime(id) ON DELETE CASCADE,
    issue_id TEXT NOT NULL REFERENCES issue(id) ON DELETE CASCADE,
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

INSERT INTO agent_task_queue_old SELECT
    id, agent_id, runtime_id, issue_id,
    status, priority, context, session_id, work_dir,
    trigger_comment_id, dispatched_at, started_at, completed_at,
    result, error, created_at
FROM agent_task_queue WHERE issue_id IS NOT NULL;

DROP TABLE agent_task_queue;
ALTER TABLE agent_task_queue_old RENAME TO agent_task_queue;

CREATE INDEX IF NOT EXISTS idx_agent_task_queue_agent ON agent_task_queue(agent_id, status);
CREATE INDEX IF NOT EXISTS idx_agent_task_queue_runtime_pending ON agent_task_queue(runtime_id, priority, created_at);
DROP INDEX IF EXISTS idx_agent_task_queue_channel;
