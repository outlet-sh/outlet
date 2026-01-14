-- +goose Up
-- Migration for retry worker support

-- Add retry_count and failed_at to campaign_sends
ALTER TABLE campaign_sends ADD COLUMN retry_count INTEGER DEFAULT 0;
ALTER TABLE campaign_sends ADD COLUMN failed_at TEXT;

-- Create index for failed sends lookup
CREATE INDEX IF NOT EXISTS idx_campaign_sends_failed ON campaign_sends(status, retry_count) WHERE status = 'failed';

-- +goose Down
-- SQLite doesn't support DROP COLUMN easily, so we leave this as a no-op
-- The columns will remain but won't be used
