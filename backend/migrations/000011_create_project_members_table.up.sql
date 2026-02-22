CREATE TABLE project_members (
    project_id INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    role       INTEGER NOT NULL DEFAULT 3,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (project_id, user_id),

    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_project_members_user ON project_members(user_id);
