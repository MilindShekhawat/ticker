-- Ticket tags junction table
CREATE TABLE IF NOT EXISTS ticket_tags (
    ticket_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,

    PRIMARY KEY (ticket_id, tag_id),
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Index for reverse lookup (find tickets by tag)
CREATE INDEX IF NOT EXISTS idx_ticket_tags_tag ON ticket_tags(tag_id);
