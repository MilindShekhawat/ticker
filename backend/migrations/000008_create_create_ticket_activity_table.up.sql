-- Ticket activity table
CREATE TABLE IF NOT EXISTS ticket_activity (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id INTEGER NOT NULL,
    actor_id INTEGER NOT NULL,
    activity_type INTEGER NOT NULL,
    metadata TEXT,  -- JSON
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES users(id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_activity_ticket ON ticket_activity(ticket_id);
CREATE INDEX IF NOT EXISTS idx_activity_actor ON ticket_activity(actor_id);
CREATE INDEX IF NOT EXISTS idx_activity_type ON ticket_activity(activity_type);
CREATE INDEX IF NOT EXISTS idx_activity_created ON ticket_activity(created_at);
