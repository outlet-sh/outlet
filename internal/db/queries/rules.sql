-- =====================================================
-- Org Rules
-- =====================================================

-- name: GetOrgRules :many
-- Get all enabled rules for an organization, ordered by salience
SELECT * FROM org_rules
WHERE org_id = sqlc.arg(org_id) AND enabled = 1
ORDER BY salience DESC, created_at;

-- name: GetOrgRulesByCategory :many
-- Get enabled rules for an organization filtered by category
SELECT * FROM org_rules
WHERE org_id = sqlc.arg(org_id) AND category = sqlc.arg(category) AND enabled = 1
ORDER BY salience DESC;

-- name: GetEntityRules :many
-- Get enabled rules for a specific entity (e.g., email_list, product)
SELECT * FROM org_rules
WHERE org_id = sqlc.arg(org_id) AND entity_type = sqlc.arg(entity_type) AND entity_id = sqlc.arg(entity_id) AND enabled = 1
ORDER BY salience DESC;

-- name: GetOrgRuleById :one
-- Get a single rule by ID
SELECT * FROM org_rules
WHERE id = sqlc.arg(id) AND org_id = sqlc.arg(org_id);

-- name: GetAllOrgRules :many
-- Get all rules for an organization (including disabled), for admin listing
SELECT * FROM org_rules
WHERE org_id = sqlc.arg(org_id)
ORDER BY category, salience DESC, name;

-- name: GetOrgRulesFiltered :many
-- Get rules with optional category and enabled filters
SELECT * FROM org_rules
WHERE org_id = sqlc.arg(org_id)
  AND (sqlc.arg(filter_category) IS NULL OR sqlc.arg(filter_category) = '' OR category = sqlc.arg(filter_category))
  AND (sqlc.arg(filter_enabled) IS NULL OR sqlc.arg(filter_enabled) = '' OR enabled = (sqlc.arg(filter_enabled) = '1'))
ORDER BY category, salience DESC, name
LIMIT sqlc.arg(page_size) OFFSET sqlc.arg(page_offset);

-- name: CountOrgRules :one
-- Count rules with optional filters
SELECT COUNT(*) as count FROM org_rules
WHERE org_id = sqlc.arg(org_id)
  AND (sqlc.arg(filter_category) IS NULL OR sqlc.arg(filter_category) = '' OR category = sqlc.arg(filter_category))
  AND (sqlc.arg(filter_enabled) IS NULL OR sqlc.arg(filter_enabled) = '' OR enabled = (sqlc.arg(filter_enabled) = '1'));

-- name: CreateOrgRule :one
-- Create a new rule
INSERT INTO org_rules (
    id, org_id, name, description, category, rule_json,
    entity_type, entity_id, enabled, salience,
    compiled_hash, validation_errors, last_validated_at,
    created_by, created_at, updated_at
)
VALUES (sqlc.arg(id), sqlc.arg(org_id), sqlc.arg(name), sqlc.arg(description), sqlc.arg(category), sqlc.arg(rule_json), sqlc.arg(entity_type), sqlc.arg(entity_id), sqlc.arg(enabled), sqlc.arg(salience), sqlc.arg(compiled_hash), sqlc.arg(validation_errors), sqlc.arg(last_validated_at), sqlc.arg(created_by), datetime('now'), datetime('now'))
RETURNING *;

-- name: UpdateOrgRule :one
-- Update an existing rule
UPDATE org_rules
SET name = sqlc.arg(name),
    description = sqlc.arg(description),
    rule_json = sqlc.arg(rule_json),
    entity_type = sqlc.arg(entity_type),
    entity_id = sqlc.arg(entity_id),
    enabled = sqlc.arg(enabled),
    salience = sqlc.arg(salience),
    compiled_hash = sqlc.arg(compiled_hash),
    validation_errors = sqlc.arg(validation_errors),
    last_validated_at = sqlc.arg(last_validated_at),
    updated_by = sqlc.arg(updated_by),
    updated_at = datetime('now')
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ToggleOrgRule :one
-- Toggle enabled state of a rule
UPDATE org_rules
SET enabled = CASE WHEN enabled = 1 THEN 0 ELSE 1 END,
    updated_by = sqlc.arg(updated_by),
    updated_at = datetime('now')
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateRuleValidation :exec
-- Update just the validation fields (after recompiling)
UPDATE org_rules
SET compiled_hash = sqlc.arg(compiled_hash),
    validation_errors = sqlc.arg(validation_errors),
    last_validated_at = datetime('now')
WHERE id = sqlc.arg(id);

-- name: DeleteOrgRule :exec
-- Delete a rule
DELETE FROM org_rules WHERE id = sqlc.arg(id) AND org_id = sqlc.arg(org_id);

-- name: GetRulesWithStaleValidation :many
-- Get rules that need revalidation (hash doesn't match content)
SELECT * FROM org_rules
WHERE org_id = sqlc.arg(org_id) AND enabled = 1
  AND (compiled_hash IS NULL OR last_validated_at < updated_at);

-- =====================================================
-- Automation Log
-- =====================================================

-- name: LogAutomation :one
-- Log a rule execution (idempotent - skips if already processed)
INSERT INTO automation_log (
    id, org_id, event_id, event_type, event_payload,
    rule_id, rule_name, rule_snapshot,
    actions_executed, success, error_message, execution_time_ms, created_at
)
VALUES (sqlc.arg(id), sqlc.arg(org_id), sqlc.arg(event_id), sqlc.arg(event_type), sqlc.arg(event_payload), sqlc.arg(rule_id), sqlc.arg(rule_name), sqlc.arg(rule_snapshot), sqlc.arg(actions_executed), sqlc.arg(success), sqlc.arg(error_message), sqlc.arg(execution_time_ms), datetime('now'))
ON CONFLICT (org_id, event_id, rule_id) DO NOTHING
RETURNING *;

-- name: CheckEventProcessed :one
-- Check if an event has already been processed by a rule
SELECT CASE WHEN COUNT(*) > 0 THEN 1 ELSE 0 END as found
FROM automation_log WHERE org_id = sqlc.arg(org_id) AND event_id = sqlc.arg(event_id) AND rule_id = sqlc.arg(rule_id);

-- name: GetAutomationLog :many
-- Get recent automation log entries for an organization
SELECT * FROM automation_log
WHERE org_id = sqlc.arg(org_id)
ORDER BY created_at DESC
LIMIT sqlc.arg(page_size) OFFSET sqlc.arg(page_offset);

-- name: GetAutomationLogFiltered :many
-- Get automation log with filters
SELECT * FROM automation_log
WHERE org_id = sqlc.arg(org_id)
  AND (sqlc.arg(filter_event_type) IS NULL OR sqlc.arg(filter_event_type) = '' OR event_type = sqlc.arg(filter_event_type))
  AND (sqlc.arg(filter_rule_id) IS NULL OR sqlc.arg(filter_rule_id) = '' OR rule_id = sqlc.arg(filter_rule_id))
  AND (sqlc.arg(filter_success) IS NULL OR sqlc.arg(filter_success) = '' OR success = (sqlc.arg(filter_success) = '1'))
ORDER BY created_at DESC
LIMIT sqlc.arg(page_size) OFFSET sqlc.arg(page_offset);

-- name: CountAutomationLog :one
-- Count automation log entries with optional filters
SELECT COUNT(*) as count FROM automation_log
WHERE org_id = sqlc.arg(org_id)
  AND (sqlc.arg(filter_event_type) IS NULL OR sqlc.arg(filter_event_type) = '' OR event_type = sqlc.arg(filter_event_type))
  AND (sqlc.arg(filter_rule_id) IS NULL OR sqlc.arg(filter_rule_id) = '' OR rule_id = sqlc.arg(filter_rule_id))
  AND (sqlc.arg(filter_success) IS NULL OR sqlc.arg(filter_success) = '' OR success = (sqlc.arg(filter_success) = '1'));

-- name: GetAutomationLogByRule :many
-- Get automation log entries for a specific rule
SELECT * FROM automation_log
WHERE rule_id = sqlc.arg(rule_id)
ORDER BY created_at DESC
LIMIT sqlc.arg(page_size) OFFSET sqlc.arg(page_offset);

-- name: GetAutomationLogStats :one
-- Get stats for automation log (for dashboard)
SELECT
    COUNT(*) as total,
    SUM(CASE WHEN success = 1 THEN 1 ELSE 0 END) as successful,
    SUM(CASE WHEN success = 0 THEN 1 ELSE 0 END) as failed,
    AVG(execution_time_ms) as avg_execution_time_ms
FROM automation_log
WHERE org_id = sqlc.arg(org_id)
  AND created_at > datetime('now', '-24 hours');

-- =====================================================
-- Rule Templates
-- =====================================================

-- name: GetRuleTemplates :many
-- Get all rule templates
SELECT * FROM rule_templates
ORDER BY is_recommended DESC, category, name;

-- name: GetRuleTemplatesByCategory :many
-- Get rule templates filtered by category
SELECT * FROM rule_templates
WHERE category = sqlc.arg(category)
ORDER BY is_recommended DESC, name;

-- name: GetRuleTemplateById :one
-- Get a single rule template by ID
SELECT * FROM rule_templates
WHERE id = sqlc.arg(id);

-- name: CreateRuleTemplate :one
-- Create a new rule template (platform admin only)
INSERT INTO rule_templates (id, name, description, category, rule_json, configurable_params, is_recommended, created_at)
VALUES (sqlc.arg(id), sqlc.arg(name), sqlc.arg(description), sqlc.arg(category), sqlc.arg(rule_json), sqlc.arg(configurable_params), sqlc.arg(is_recommended), datetime('now'))
RETURNING *;

-- name: UpdateRuleTemplate :one
-- Update a rule template
UPDATE rule_templates
SET name = sqlc.arg(name),
    description = sqlc.arg(description),
    category = sqlc.arg(category),
    rule_json = sqlc.arg(rule_json),
    configurable_params = sqlc.arg(configurable_params),
    is_recommended = sqlc.arg(is_recommended)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteRuleTemplate :exec
-- Delete a rule template
DELETE FROM rule_templates WHERE id = sqlc.arg(id);

-- name: GetDefaultRuleTemplates :many
-- Get all rule templates marked as defaults
SELECT * FROM rule_templates
WHERE is_default = 1
ORDER BY category, name;
