-- up
CREATE TABLE IF NOT EXISTS projects (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    key_prefix  TEXT    NOT NULL,
    name        TEXT    NOT NULL,
    description TEXT,
    created_by  INTEGER NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id),
    UNIQUE (key_prefix)
);

-- Used in: ListByUser, GetByUser → JOIN project_members WHERE pm.user_id = ?
CREATE INDEX IF NOT EXISTS idx_project_members_user_id    ON project_members(user_id);

-- Used in: ListByUser, GetByUser, UpdateByUser, DeleteByUser → JOIN/subquery on pm.project_id
CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members(project_id);
