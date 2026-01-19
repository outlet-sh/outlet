-- +goose Up
-- Domain identities for AWS SES verification

CREATE TABLE IF NOT EXISTS domain_identities (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    domain TEXT NOT NULL,
    verification_status TEXT DEFAULT 'pending' CHECK (verification_status IN ('pending', 'success', 'failed', 'temporary_failure', 'not_started')),
    dkim_status TEXT DEFAULT 'pending' CHECK (dkim_status IN ('pending', 'success', 'failed', 'temporary_failure', 'not_started')),
    verification_token TEXT,
    dkim_tokens TEXT, -- JSON array of DKIM tokens
    dns_records TEXT, -- JSON array of all DNS records needed
    mail_from_domain TEXT,
    mail_from_status TEXT DEFAULT 'not_started' CHECK (mail_from_status IN ('pending', 'success', 'failed', 'not_started')),
    last_checked_at TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    UNIQUE(org_id, domain)
);

CREATE INDEX IF NOT EXISTS idx_domain_identities_org_id ON domain_identities(org_id);
CREATE INDEX IF NOT EXISTS idx_domain_identities_domain ON domain_identities(domain);

-- +goose Down
DROP TABLE IF EXISTS domain_identities;
