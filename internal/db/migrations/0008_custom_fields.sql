-- +goose Up
-- Custom fields allow lists to define additional subscriber data fields
-- These can be used for personalization in emails via merge tags like {{first_name}}

-- Custom field definitions per list
CREATE TABLE IF NOT EXISTS custom_fields (
    id TEXT PRIMARY KEY,
    list_id INTEGER NOT NULL REFERENCES email_lists(id) ON DELETE CASCADE,
    name TEXT NOT NULL,           -- Display name (e.g., "First Name")
    field_key TEXT NOT NULL,      -- Merge tag key (e.g., "first_name") - lowercase, underscores
    field_type TEXT NOT NULL DEFAULT 'text',  -- text, number, date, dropdown
    options TEXT,                 -- JSON array for dropdown options (e.g., '["Option 1", "Option 2"]')
    required INTEGER NOT NULL DEFAULT 0,
    default_value TEXT,
    placeholder TEXT,             -- Placeholder text for form input
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    UNIQUE(list_id, field_key)
);

CREATE INDEX IF NOT EXISTS idx_custom_fields_list_id ON custom_fields(list_id);

-- Custom field values per list subscriber
CREATE TABLE IF NOT EXISTS custom_field_values (
    id TEXT PRIMARY KEY,
    subscriber_id TEXT NOT NULL REFERENCES list_subscribers(id) ON DELETE CASCADE,
    field_id TEXT NOT NULL REFERENCES custom_fields(id) ON DELETE CASCADE,
    value TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    UNIQUE(subscriber_id, field_id)
);

CREATE INDEX IF NOT EXISTS idx_custom_field_values_subscriber_id ON custom_field_values(subscriber_id);
CREATE INDEX IF NOT EXISTS idx_custom_field_values_field_id ON custom_field_values(field_id);

-- +goose Down
DROP INDEX IF EXISTS idx_custom_field_values_field_id;
DROP INDEX IF EXISTS idx_custom_field_values_subscriber_id;
DROP TABLE IF EXISTS custom_field_values;
DROP INDEX IF EXISTS idx_custom_fields_list_id;
DROP TABLE IF EXISTS custom_fields;
