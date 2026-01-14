-- +goose Up
-- Outlet.sh - Blocklist Enhancements for Newsletter Feature Parity

-- =============================================================================
-- BLOCKED DOMAINS (Block entire domains like disposable email providers)
-- =============================================================================

CREATE TABLE IF NOT EXISTS blocked_domains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    domain TEXT NOT NULL COLLATE NOCASE,
    reason TEXT,
    block_attempts INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    UNIQUE(org_id, domain)
);

CREATE INDEX IF NOT EXISTS idx_blocked_domains_org ON blocked_domains(org_id);
CREATE INDEX IF NOT EXISTS idx_blocked_domains_domain ON blocked_domains(domain);

-- =============================================================================
-- SUPPRESSION LIST (Manual email suppression - separate from bounces)
-- =============================================================================

CREATE TABLE IF NOT EXISTS suppression_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    email_lower TEXT NOT NULL COLLATE NOCASE,
    reason TEXT,
    source TEXT DEFAULT 'manual', -- manual, import, api
    block_attempts INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now')),
    UNIQUE(org_id, email_lower)
);

CREATE INDEX IF NOT EXISTS idx_suppression_list_org ON suppression_list(org_id);
CREATE INDEX IF NOT EXISTS idx_suppression_list_email ON suppression_list(email_lower);

-- =============================================================================
-- IMPORT JOBS (Track CSV import progress)
-- =============================================================================

CREATE TABLE IF NOT EXISTS import_jobs (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    list_id INTEGER REFERENCES email_lists(id) ON DELETE SET NULL,
    type TEXT NOT NULL CHECK (type IN ('subscribers', 'suppression', 'blocked_domains')),
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    filename TEXT NOT NULL,
    total_rows INTEGER DEFAULT 0,
    processed_rows INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    error_count INTEGER DEFAULT 0,
    skip_count INTEGER DEFAULT 0,
    errors TEXT, -- JSON array of error messages
    options TEXT, -- JSON import options (update_existing, etc.)
    started_at TEXT,
    completed_at TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_import_jobs_org ON import_jobs(org_id);
CREATE INDEX IF NOT EXISTS idx_import_jobs_status ON import_jobs(status);

-- +goose Down
DROP TABLE IF EXISTS import_jobs;
DROP TABLE IF EXISTS suppression_list;
DROP TABLE IF EXISTS blocked_domains;
