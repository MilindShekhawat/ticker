-- Tickets table
-- Based on updated schema with project relationship and soft deletes
CREATE TABLE IF NOT EXISTS tickets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    ticket_number INTEGER NOT NULL,

    title TEXT NOT NULL,
    description TEXT,

    status_id INTEGER NOT NULL,
    priority_id INTEGER NOT NULL,
    position INTEGER DEFAULT 0,

    assignee_id INTEGER,
    created_by INTEGER NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,

    -- Foreign keys
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (status_id) REFERENCES statuses(id),
    FOREIGN KEY (priority_id) REFERENCES priorities(id),
    FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id),

    -- Ensure unique ticket numbers per project
    UNIQUE(project_id, ticket_number)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_tickets_project ON tickets(project_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status_id);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets(priority_id);
CREATE INDEX IF NOT EXISTS idx_tickets_assignee ON tickets(assignee_id);
CREATE INDEX IF NOT EXISTS idx_tickets_creator ON tickets(created_by);
CREATE INDEX IF NOT EXISTS idx_tickets_created_at ON tickets(created_at);
CREATE INDEX IF NOT EXISTS idx_tickets_deleted_at ON tickets(deleted_at);

-- Composite index for common query pattern (project + status)
CREATE INDEX IF NOT EXISTS idx_tickets_project_status ON tickets(project_id, status_id);
