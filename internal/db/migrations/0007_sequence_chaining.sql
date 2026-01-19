-- +goose Up
-- Add sequence chaining support
-- Allows a sequence to automatically enroll contacts into another sequence on completion

ALTER TABLE email_sequences ADD COLUMN on_completion_sequence_id TEXT REFERENCES email_sequences(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_email_sequences_on_completion ON email_sequences(on_completion_sequence_id);

-- +goose Down
DROP INDEX IF EXISTS idx_email_sequences_on_completion;
-- SQLite doesn't support DROP COLUMN, so we leave the column in place for down migration
