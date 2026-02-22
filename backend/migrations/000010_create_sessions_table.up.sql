CREATE TABLE sessions (
    id            TEXT PRIMARY KEY,
    user_id       INTEGER NOT NULL,
    ip_address    TEXT,
    user_agent    TEXT,
    expires_at    DATETIME NOT NULL,
    revoked_at    DATETIME NULL,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for performance

-- Fast lookup by user (logout all sessions, admin actions)
CREATE INDEX idx_sessions_user_id ON sessions(user_id);

-- Efficient expiry cleanup
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Only active sessions lookup optimization
CREATE INDEX idx_sessions_active ON sessions(id, revoked_at, expires_at);
