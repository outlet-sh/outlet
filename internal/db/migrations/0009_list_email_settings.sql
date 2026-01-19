-- +goose Up
-- Extended list email settings for thank you emails, goodbye emails, and unsubscribe behavior

-- Thank you email settings (sent after subscription)
ALTER TABLE email_lists ADD COLUMN thank_you_email_enabled INTEGER DEFAULT 0;
ALTER TABLE email_lists ADD COLUMN thank_you_email_subject TEXT;
ALTER TABLE email_lists ADD COLUMN thank_you_email_body TEXT;

-- Already subscribed redirect URL
ALTER TABLE email_lists ADD COLUMN already_subscribed_url TEXT;

-- Goodbye email settings (sent after unsubscription)
ALTER TABLE email_lists ADD COLUMN goodbye_email_enabled INTEGER DEFAULT 0;
ALTER TABLE email_lists ADD COLUMN goodbye_email_subject TEXT;
ALTER TABLE email_lists ADD COLUMN goodbye_email_body TEXT;

-- Unsubscribe behavior settings
-- 'single' = immediate unsubscribe, 'double' = require confirmation click
ALTER TABLE email_lists ADD COLUMN unsubscribe_behavior TEXT DEFAULT 'single';
-- 'list' = unsubscribe from this list only, 'all' = unsubscribe from all lists
ALTER TABLE email_lists ADD COLUMN unsubscribe_scope TEXT DEFAULT 'list';

-- +goose Down
ALTER TABLE email_lists DROP COLUMN unsubscribe_scope;
ALTER TABLE email_lists DROP COLUMN unsubscribe_behavior;
ALTER TABLE email_lists DROP COLUMN goodbye_email_body;
ALTER TABLE email_lists DROP COLUMN goodbye_email_subject;
ALTER TABLE email_lists DROP COLUMN goodbye_email_enabled;
ALTER TABLE email_lists DROP COLUMN already_subscribed_url;
ALTER TABLE email_lists DROP COLUMN thank_you_email_body;
ALTER TABLE email_lists DROP COLUMN thank_you_email_subject;
ALTER TABLE email_lists DROP COLUMN thank_you_email_enabled;
