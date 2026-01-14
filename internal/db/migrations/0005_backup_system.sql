-- +goose Up
-- Migration for backup system

-- Backup history table
CREATE TABLE IF NOT EXISTS backup_history (
    id TEXT PRIMARY KEY,
    filename TEXT NOT NULL,
    file_path TEXT,
    file_size INTEGER NOT NULL DEFAULT 0,
    backup_type TEXT NOT NULL CHECK (backup_type IN ('manual', 'scheduled', 'cli')),
    storage_type TEXT NOT NULL CHECK (storage_type IN ('local', 's3')),
    s3_bucket TEXT,
    s3_key TEXT,
    status TEXT NOT NULL CHECK (status IN ('pending', 'in_progress', 'completed', 'failed')),
    error_message TEXT,
    created_by TEXT REFERENCES users(id),
    started_at TEXT DEFAULT (datetime('now')),
    completed_at TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_backup_history_status ON backup_history(status);
CREATE INDEX IF NOT EXISTS idx_backup_history_backup_type ON backup_history(backup_type);
CREATE INDEX IF NOT EXISTS idx_backup_history_created_at ON backup_history(created_at);

-- Backup settings in platform_settings (using existing table)
-- Categories: backup_local, backup_s3, backup_schedule

-- +goose Down
DROP TABLE IF EXISTS backup_history;
