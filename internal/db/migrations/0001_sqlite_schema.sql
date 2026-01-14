-- +goose Up
-- Outlet.sh - Self-hosted Email Newsletter Platform (SQLite Schema)
-- Converted from PostgreSQL for single-binary deployment

-- =============================================================================
-- CORE AUTH TABLES
-- =============================================================================

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE COLLATE NOCASE,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'manager', 'agent', 'super_admin')),
    status TEXT DEFAULT 'pending' NOT NULL CHECK (status IN ('active', 'inactive', 'pending')),
    email_verified INTEGER DEFAULT 0 NOT NULL,
    phone TEXT,
    avatar_url TEXT,
    last_login_at TEXT,
    password_changed_at TEXT DEFAULT (datetime('now')),
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

CREATE TABLE IF NOT EXISTS auth_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    token_type TEXT NOT NULL CHECK (token_type IN ('password_reset', 'email_verification', 'invitation')),
    expires_at TEXT NOT NULL,
    used_at TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_auth_tokens_token ON auth_tokens(token);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_user_id ON auth_tokens(user_id);

-- =============================================================================
-- ORGANIZATIONS (Multi-tenancy - like Sendy "Brands")
-- =============================================================================

CREATE TABLE IF NOT EXISTS organizations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    api_key TEXT UNIQUE NOT NULL,
    api_key_created_at TEXT DEFAULT (datetime('now')),
    from_name TEXT,
    from_email TEXT,
    reply_to TEXT,
    max_contacts INTEGER DEFAULT 1000,
    settings TEXT DEFAULT '{}',
    app_url TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_organizations_api_key ON organizations(api_key);
CREATE INDEX IF NOT EXISTS idx_organizations_slug ON organizations(slug);

CREATE TABLE IF NOT EXISTS org_users (
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT DEFAULT 'member',
    created_at TEXT DEFAULT (datetime('now')),
    PRIMARY KEY (org_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_org_users_user_id ON org_users(user_id);

CREATE TABLE IF NOT EXISTS org_invitations (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    role TEXT DEFAULT 'member',
    token TEXT NOT NULL UNIQUE,
    expires_at TEXT NOT NULL,
    accepted_at TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_org_invitations_token ON org_invitations(token);
CREATE INDEX IF NOT EXISTS idx_org_invitations_org_id ON org_invitations(org_id);

-- =============================================================================
-- MCP INTEGRATION (Model Context Protocol)
-- =============================================================================

CREATE TABLE IF NOT EXISTS mcp_api_keys (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    key_hash TEXT NOT NULL UNIQUE,
    key_prefix TEXT NOT NULL,
    scopes TEXT DEFAULT '["mcp:full"]',
    last_used_at TEXT,
    expires_at TEXT,
    revoked_at TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_mcp_api_keys_user_id ON mcp_api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_mcp_api_keys_key_prefix ON mcp_api_keys(key_prefix);

CREATE TABLE IF NOT EXISTS mcp_oauth_clients (
    id TEXT PRIMARY KEY,
    client_id TEXT NOT NULL UNIQUE,
    client_secret_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    redirect_uris TEXT NOT NULL,
    scopes TEXT DEFAULT '["mcp:full"]',
    is_confidential INTEGER DEFAULT 1,
    is_active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_mcp_oauth_clients_client_id ON mcp_oauth_clients(client_id);

CREATE TABLE IF NOT EXISTS mcp_oauth_codes (
    id TEXT PRIMARY KEY,
    client_id TEXT NOT NULL REFERENCES mcp_oauth_clients(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash TEXT NOT NULL UNIQUE,
    redirect_uri TEXT NOT NULL,
    scopes TEXT NOT NULL,
    code_challenge TEXT,
    code_challenge_method TEXT,
    expires_at TEXT NOT NULL,
    used_at TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_mcp_oauth_codes_client_id ON mcp_oauth_codes(client_id);
CREATE INDEX IF NOT EXISTS idx_mcp_oauth_codes_user_id ON mcp_oauth_codes(user_id);

CREATE TABLE IF NOT EXISTS mcp_oauth_tokens (
    id TEXT PRIMARY KEY,
    client_id TEXT NOT NULL REFERENCES mcp_oauth_clients(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_token_hash TEXT NOT NULL UNIQUE,
    refresh_token_hash TEXT UNIQUE,
    scopes TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    refresh_expires_at TEXT,
    revoked_at TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_mcp_oauth_tokens_client_id ON mcp_oauth_tokens(client_id);
CREATE INDEX IF NOT EXISTS idx_mcp_oauth_tokens_user_id ON mcp_oauth_tokens(user_id);

CREATE TABLE IF NOT EXISTS mcp_sessions (
    session_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    org_id TEXT REFERENCES organizations(id) ON DELETE SET NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_mcp_sessions_user_id ON mcp_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_mcp_sessions_updated_at ON mcp_sessions(updated_at);

-- =============================================================================
-- CONTACTS & SUBSCRIBERS
-- =============================================================================

CREATE TABLE IF NOT EXISTS contacts (
    id TEXT PRIMARY KEY,
    org_id TEXT REFERENCES organizations(id),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    source TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    email_verified INTEGER DEFAULT 0 NOT NULL,
    verification_token TEXT,
    verification_sent_at TEXT,
    verified_at TEXT,
    unsubscribed_at TEXT,
    blocked_at TEXT,
    status TEXT DEFAULT 'new',
    gdpr_consent INTEGER DEFAULT 0,
    gdpr_consent_at TEXT
);

CREATE INDEX IF NOT EXISTS idx_contacts_org_id ON contacts(org_id);
CREATE INDEX IF NOT EXISTS idx_contacts_email ON contacts(email);
CREATE INDEX IF NOT EXISTS idx_contacts_created_at ON contacts(created_at);
CREATE INDEX IF NOT EXISTS idx_contacts_status ON contacts(status);
CREATE INDEX IF NOT EXISTS idx_contacts_verification_token ON contacts(verification_token);

CREATE TABLE IF NOT EXISTS contact_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    contact_id TEXT REFERENCES contacts(id) ON DELETE CASCADE,
    tag TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    UNIQUE(contact_id, tag)
);

CREATE INDEX IF NOT EXISTS idx_contact_tags_contact ON contact_tags(contact_id);
CREATE INDEX IF NOT EXISTS idx_contact_tags_tag ON contact_tags(tag);

-- =============================================================================
-- EMAIL LISTS & SUBSCRIBERS
-- =============================================================================

CREATE TABLE IF NOT EXISTS email_lists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    double_optin INTEGER DEFAULT 1,
    confirmation_email_subject TEXT,
    confirmation_email_body TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    UNIQUE(org_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_email_lists_org_id ON email_lists(org_id);

CREATE TABLE IF NOT EXISTS list_subscribers (
    id TEXT PRIMARY KEY,
    list_id INTEGER NOT NULL REFERENCES email_lists(id) ON DELETE CASCADE,
    contact_id TEXT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    status TEXT DEFAULT 'active',
    verification_token TEXT,
    verification_sent_at TEXT,
    verified_at TEXT,
    subscribed_at TEXT DEFAULT (datetime('now')),
    unsubscribed_at TEXT,
    UNIQUE(list_id, contact_id)
);

CREATE INDEX IF NOT EXISTS idx_list_subscribers_list_id ON list_subscribers(list_id);
CREATE INDEX IF NOT EXISTS idx_list_subscribers_contact_id ON list_subscribers(contact_id);
CREATE INDEX IF NOT EXISTS idx_list_subscribers_verification_token ON list_subscribers(verification_token);

CREATE TABLE IF NOT EXISTS subscriber_custom_fields (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    list_id INTEGER NOT NULL REFERENCES email_lists(id) ON DELETE CASCADE,
    contact_id TEXT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    field_name TEXT NOT NULL,
    field_value TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    UNIQUE(list_id, contact_id, field_name)
);

CREATE INDEX IF NOT EXISTS idx_subscriber_custom_fields_list ON subscriber_custom_fields(list_id);
CREATE INDEX IF NOT EXISTS idx_subscriber_custom_fields_contact ON subscriber_custom_fields(contact_id);

-- =============================================================================
-- EMAIL DESIGNS (Reusable Templates)
-- =============================================================================

CREATE TABLE IF NOT EXISTS email_designs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    category TEXT DEFAULT 'general',
    html_body TEXT NOT NULL,
    plain_text TEXT,
    thumbnail_url TEXT,
    is_active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    UNIQUE(org_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_email_designs_org ON email_designs(org_id);
CREATE INDEX IF NOT EXISTS idx_email_designs_category ON email_designs(org_id, category);

-- =============================================================================
-- EMAIL SEQUENCES (Autoresponders)
-- =============================================================================

CREATE TABLE IF NOT EXISTS email_sequences (
    id TEXT PRIMARY KEY,
    org_id TEXT REFERENCES organizations(id),
    list_id INTEGER REFERENCES email_lists(id) ON DELETE SET NULL,
    slug TEXT NOT NULL,
    name TEXT NOT NULL,
    trigger_event TEXT NOT NULL,
    is_active INTEGER DEFAULT 1,
    send_hour INTEGER DEFAULT 9,
    send_timezone TEXT DEFAULT 'America/New_York',
    sequence_type TEXT DEFAULT 'lifecycle',
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_email_sequences_org_id ON email_sequences(org_id);
CREATE INDEX IF NOT EXISTS idx_email_sequences_list_id ON email_sequences(list_id);

CREATE TABLE IF NOT EXISTS email_templates (
    id TEXT PRIMARY KEY,
    org_id TEXT REFERENCES organizations(id),
    sequence_id TEXT REFERENCES email_sequences(id) ON DELETE CASCADE,
    design_id INTEGER REFERENCES email_designs(id) ON DELETE SET NULL,
    position INTEGER NOT NULL,
    delay_hours INTEGER DEFAULT 0 NOT NULL,
    subject TEXT NOT NULL,
    html_body TEXT NOT NULL,
    plain_text TEXT,
    template_type TEXT DEFAULT 'email',
    is_active INTEGER DEFAULT 1,
    is_transactional INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now')),
    UNIQUE(sequence_id, position)
);

CREATE INDEX IF NOT EXISTS idx_email_templates_org_id ON email_templates(org_id);
CREATE INDEX IF NOT EXISTS idx_email_templates_sequence_id ON email_templates(sequence_id);

CREATE TABLE IF NOT EXISTS sequence_entry_rules (
    id TEXT PRIMARY KEY,
    sequence_id TEXT NOT NULL REFERENCES email_sequences(id) ON DELETE CASCADE,
    trigger_type TEXT NOT NULL,
    source_id TEXT NOT NULL,
    priority INTEGER DEFAULT 0,
    is_active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_sequence_entry_rules_sequence ON sequence_entry_rules(sequence_id);
CREATE INDEX IF NOT EXISTS idx_sequence_entry_rules_source ON sequence_entry_rules(source_id);
CREATE INDEX IF NOT EXISTS idx_sequence_entry_rules_trigger ON sequence_entry_rules(trigger_type);

CREATE TABLE IF NOT EXISTS contact_sequence_state (
    id TEXT PRIMARY KEY,
    contact_id TEXT REFERENCES contacts(id) ON DELETE CASCADE,
    sequence_id TEXT REFERENCES email_sequences(id) ON DELETE CASCADE,
    current_position INTEGER DEFAULT 0,
    is_active INTEGER DEFAULT 1,
    started_at TEXT DEFAULT (datetime('now')),
    completed_at TEXT,
    unsubscribed_at TEXT,
    paused_at TEXT,
    UNIQUE(contact_id, sequence_id)
);

CREATE INDEX IF NOT EXISTS idx_contact_sequence_state_contact ON contact_sequence_state(contact_id);

CREATE TABLE IF NOT EXISTS email_queue (
    id TEXT PRIMARY KEY,
    contact_id TEXT REFERENCES contacts(id) ON DELETE CASCADE,
    template_id TEXT REFERENCES email_templates(id) ON DELETE CASCADE,
    scheduled_for TEXT NOT NULL,
    sent_at TEXT,
    status TEXT DEFAULT 'pending',
    error_message TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    tracking_token TEXT UNIQUE,
    opened_at TEXT,
    open_count INTEGER DEFAULT 0 NOT NULL,
    clicked_at TEXT,
    click_count INTEGER DEFAULT 0 NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_email_queue_contact ON email_queue(contact_id);
CREATE INDEX IF NOT EXISTS idx_email_queue_scheduled ON email_queue(status, scheduled_for);
CREATE INDEX IF NOT EXISTS idx_email_queue_tracking_token ON email_queue(tracking_token);

CREATE TABLE IF NOT EXISTS email_clicks (
    id TEXT PRIMARY KEY,
    email_queue_id TEXT REFERENCES email_queue(id) ON DELETE CASCADE,
    contact_id TEXT REFERENCES contacts(id) ON DELETE CASCADE,
    link_url TEXT NOT NULL,
    link_name TEXT,
    clicked_at TEXT DEFAULT (datetime('now')) NOT NULL,
    user_agent TEXT,
    ip_address TEXT
);

CREATE INDEX IF NOT EXISTS idx_email_clicks_queue ON email_clicks(email_queue_id);
CREATE INDEX IF NOT EXISTS idx_email_clicks_contact ON email_clicks(contact_id);

-- =============================================================================
-- EMAIL CAMPAIGNS (Broadcasts)
-- =============================================================================

CREATE TABLE IF NOT EXISTS email_campaigns (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    design_id INTEGER REFERENCES email_designs(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    subject TEXT NOT NULL,
    preview_text TEXT,
    from_name TEXT,
    from_email TEXT,
    reply_to TEXT,
    html_body TEXT NOT NULL,
    plain_text TEXT,
    list_ids TEXT DEFAULT '[]',
    exclude_list_ids TEXT DEFAULT '[]',
    segment_filter TEXT,
    status TEXT DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'sending', 'sent', 'paused', 'cancelled')),
    scheduled_at TEXT,
    started_at TEXT,
    completed_at TEXT,
    track_opens INTEGER DEFAULT 1,
    track_clicks INTEGER DEFAULT 1,
    recipients_count INTEGER DEFAULT 0,
    sent_count INTEGER DEFAULT 0,
    delivered_count INTEGER DEFAULT 0,
    opened_count INTEGER DEFAULT 0,
    clicked_count INTEGER DEFAULT 0,
    bounced_count INTEGER DEFAULT 0,
    complained_count INTEGER DEFAULT 0,
    unsubscribed_count INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_campaigns_org ON email_campaigns(org_id);
CREATE INDEX IF NOT EXISTS idx_campaigns_status ON email_campaigns(org_id, status);

CREATE TABLE IF NOT EXISTS campaign_sends (
    id TEXT PRIMARY KEY,
    campaign_id TEXT NOT NULL REFERENCES email_campaigns(id) ON DELETE CASCADE,
    contact_id TEXT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    list_id INTEGER REFERENCES email_lists(id) ON DELETE SET NULL,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'delivered', 'bounced', 'failed')),
    sent_at TEXT,
    delivered_at TEXT,
    tracking_token TEXT UNIQUE,
    opened_at TEXT,
    open_count INTEGER DEFAULT 0,
    clicked_at TEXT,
    click_count INTEGER DEFAULT 0,
    error_message TEXT,
    bounce_type TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    UNIQUE(campaign_id, contact_id)
);

CREATE INDEX IF NOT EXISTS idx_campaign_sends_campaign ON campaign_sends(campaign_id);
CREATE INDEX IF NOT EXISTS idx_campaign_sends_contact ON campaign_sends(contact_id);
CREATE INDEX IF NOT EXISTS idx_campaign_sends_tracking ON campaign_sends(tracking_token);

CREATE TABLE IF NOT EXISTS campaign_clicks (
    id TEXT PRIMARY KEY,
    campaign_send_id TEXT NOT NULL REFERENCES campaign_sends(id) ON DELETE CASCADE,
    campaign_id TEXT NOT NULL REFERENCES email_campaigns(id) ON DELETE CASCADE,
    contact_id TEXT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    link_url TEXT NOT NULL,
    link_name TEXT,
    clicked_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_campaign_clicks_send ON campaign_clicks(campaign_send_id);
CREATE INDEX IF NOT EXISTS idx_campaign_clicks_campaign ON campaign_clicks(campaign_id);

-- =============================================================================
-- TRANSACTIONAL EMAILS
-- =============================================================================

CREATE TABLE IF NOT EXISTS transactional_emails (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    design_id INTEGER REFERENCES email_designs(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT,
    subject TEXT NOT NULL,
    html_body TEXT NOT NULL,
    plain_text TEXT,
    from_name TEXT,
    from_email TEXT,
    reply_to TEXT,
    is_active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    UNIQUE(org_id, slug)
);

CREATE INDEX IF NOT EXISTS idx_transactional_org ON transactional_emails(org_id);
CREATE INDEX IF NOT EXISTS idx_transactional_slug ON transactional_emails(org_id, slug);

CREATE TABLE IF NOT EXISTS transactional_sends (
    id TEXT PRIMARY KEY,
    template_id TEXT NOT NULL REFERENCES transactional_emails(id) ON DELETE CASCADE,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    to_email TEXT NOT NULL,
    to_name TEXT,
    contact_id TEXT REFERENCES contacts(id) ON DELETE SET NULL,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'delivered', 'bounced', 'failed')),
    sent_at TEXT,
    delivered_at TEXT,
    tracking_token TEXT UNIQUE,
    opened_at TEXT,
    clicked_at TEXT,
    context_data TEXT,
    error_message TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_transactional_sends_template ON transactional_sends(template_id);
CREATE INDEX IF NOT EXISTS idx_transactional_sends_org ON transactional_sends(org_id);
CREATE INDEX IF NOT EXISTS idx_transactional_sends_tracking ON transactional_sends(tracking_token);

-- =============================================================================
-- EMAIL DELIVERABILITY (Bounces & Complaints)
-- =============================================================================

CREATE TABLE IF NOT EXISTS email_bounce (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL,
    email_lower TEXT NOT NULL UNIQUE,
    bounce_type TEXT NOT NULL,
    bounce_subtype TEXT,
    diagnostic_code TEXT,
    source_email TEXT,
    message_id TEXT,
    raw_notification TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_email_bounce_email ON email_bounce(email);
CREATE INDEX IF NOT EXISTS idx_email_bounce_email_lower ON email_bounce(email_lower);

CREATE TABLE IF NOT EXISTS email_complaint (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL,
    email_lower TEXT NOT NULL UNIQUE,
    complaint_type TEXT,
    feedback_id TEXT,
    source_email TEXT,
    message_id TEXT,
    raw_notification TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_email_complaint_email ON email_complaint(email);
CREATE INDEX IF NOT EXISTS idx_email_complaint_email_lower ON email_complaint(email_lower);

-- =============================================================================
-- WEBHOOKS
-- =============================================================================

CREATE TABLE IF NOT EXISTS webhooks (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    secret TEXT NOT NULL,
    events TEXT NOT NULL DEFAULT '[]',
    active INTEGER DEFAULT 1,
    deliveries_total INTEGER DEFAULT 0,
    deliveries_success INTEGER DEFAULT 0,
    deliveries_failed INTEGER DEFAULT 0,
    last_delivery_at TEXT,
    last_status INTEGER,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_webhooks_org_id ON webhooks(org_id);

CREATE TABLE IF NOT EXISTS webhook_logs (
    id TEXT PRIMARY KEY,
    webhook_id TEXT NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event TEXT NOT NULL,
    payload TEXT NOT NULL,
    status_code INTEGER,
    response TEXT,
    error TEXT,
    duration_ms INTEGER,
    delivered_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_webhook_logs_webhook_id ON webhook_logs(webhook_id);
CREATE INDEX IF NOT EXISTS idx_webhook_logs_delivered_at ON webhook_logs(delivered_at);

-- =============================================================================
-- AUTOMATION RULES
-- =============================================================================

CREATE TABLE IF NOT EXISTS org_rules (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    category TEXT NOT NULL,
    rule_json TEXT NOT NULL,
    entity_type TEXT,
    entity_id TEXT,
    enabled INTEGER DEFAULT 1,
    salience INTEGER DEFAULT 0,
    compiled_hash TEXT,
    validation_errors TEXT,
    last_validated_at TEXT,
    created_by TEXT REFERENCES users(id),
    updated_by TEXT REFERENCES users(id),
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    UNIQUE(org_id, name)
);

CREATE INDEX IF NOT EXISTS idx_org_rules_org_category ON org_rules(org_id, category);

CREATE TABLE IF NOT EXISTS automation_log (
    id TEXT PRIMARY KEY,
    org_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    event_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    event_payload TEXT,
    rule_id TEXT REFERENCES org_rules(id) ON DELETE SET NULL,
    rule_name TEXT NOT NULL,
    rule_snapshot TEXT,
    actions_executed TEXT,
    success INTEGER NOT NULL,
    error_message TEXT,
    execution_time_ms INTEGER,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_automation_log_dedupe ON automation_log(org_id, event_id, rule_id);
CREATE INDEX IF NOT EXISTS idx_automation_log_org ON automation_log(org_id, created_at);

CREATE TABLE IF NOT EXISTS rule_templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    category TEXT NOT NULL,
    rule_json TEXT NOT NULL,
    configurable_params TEXT,
    is_recommended INTEGER DEFAULT 0,
    is_default INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now'))
);

-- =============================================================================
-- PLATFORM SETTINGS
-- =============================================================================

CREATE TABLE IF NOT EXISTS platform_settings (
    key TEXT PRIMARY KEY,
    value_encrypted TEXT,
    value_text TEXT,
    description TEXT,
    category TEXT NOT NULL DEFAULT 'general',
    is_sensitive INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_platform_settings_category ON platform_settings(category);

-- +goose Down
DROP TABLE IF EXISTS platform_settings;
DROP TABLE IF EXISTS rule_templates;
DROP TABLE IF EXISTS automation_log;
DROP TABLE IF EXISTS org_rules;
DROP TABLE IF EXISTS webhook_logs;
DROP TABLE IF EXISTS webhooks;
DROP TABLE IF EXISTS email_complaint;
DROP TABLE IF EXISTS email_bounce;
DROP TABLE IF EXISTS transactional_sends;
DROP TABLE IF EXISTS transactional_emails;
DROP TABLE IF EXISTS campaign_clicks;
DROP TABLE IF EXISTS campaign_sends;
DROP TABLE IF EXISTS email_campaigns;
DROP TABLE IF EXISTS email_clicks;
DROP TABLE IF EXISTS email_queue;
DROP TABLE IF EXISTS contact_sequence_state;
DROP TABLE IF EXISTS sequence_entry_rules;
DROP TABLE IF EXISTS email_templates;
DROP TABLE IF EXISTS email_sequences;
DROP TABLE IF EXISTS email_designs;
DROP TABLE IF EXISTS subscriber_custom_fields;
DROP TABLE IF EXISTS list_subscribers;
DROP TABLE IF EXISTS email_lists;
DROP TABLE IF EXISTS contact_tags;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS mcp_sessions;
DROP TABLE IF EXISTS mcp_oauth_tokens;
DROP TABLE IF EXISTS mcp_oauth_codes;
DROP TABLE IF EXISTS mcp_oauth_clients;
DROP TABLE IF EXISTS mcp_api_keys;
DROP TABLE IF EXISTS org_invitations;
DROP TABLE IF EXISTS org_users;
DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS auth_tokens;
DROP TABLE IF EXISTS users;
