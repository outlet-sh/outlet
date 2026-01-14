-- +goose Up
-- Public page settings on email_lists
ALTER TABLE email_lists ADD COLUMN public_page_enabled INTEGER DEFAULT 1;
ALTER TABLE email_lists ADD COLUMN thank_you_url TEXT;
ALTER TABLE email_lists ADD COLUMN confirm_redirect_url TEXT;
ALTER TABLE email_lists ADD COLUMN unsubscribe_redirect_url TEXT;

-- Index for list slug lookups (public pages)
CREATE INDEX IF NOT EXISTS idx_email_lists_slug ON email_lists(slug);

-- +goose Down
DROP INDEX IF EXISTS idx_email_lists_slug;
ALTER TABLE email_lists DROP COLUMN unsubscribe_redirect_url;
ALTER TABLE email_lists DROP COLUMN confirm_redirect_url;
ALTER TABLE email_lists DROP COLUMN thank_you_url;
ALTER TABLE email_lists DROP COLUMN public_page_enabled;
