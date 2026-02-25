-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key_prefix TEXT NOT NULL UNIQUE,

    name TEXT NOT NULL,
    description TEXT,

    created_by INTEGER NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Foreign keys
    FOREIGN KEY (created_by) REFERENCES users(id),

    -- Ensure unique key prefix
    UNIQUE(key_prefix)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_projects_creator ON projects(created_by);
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at);
