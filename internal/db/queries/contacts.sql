-- name: CreateContact :one
INSERT INTO contacts (
    id, org_id, name, email, source, status, created_at, updated_at
) VALUES (
    sqlc.arg(id), sqlc.arg(org_id), sqlc.arg(name), sqlc.arg(email),
    sqlc.arg(source), COALESCE(sqlc.arg(status), 'new'),
    datetime('now'), datetime('now')
) RETURNING *;

-- name: GetContactByOrgAndEmail :one
SELECT * FROM contacts
WHERE org_id = sqlc.arg(org_id) AND email = sqlc.arg(email)
ORDER BY created_at DESC LIMIT 1;

-- name: GetContact :one
SELECT * FROM contacts WHERE id = sqlc.arg(id);

-- name: GetContactByID :one
SELECT * FROM contacts WHERE id = sqlc.arg(id);

-- name: GetContactByEmail :one
SELECT * FROM contacts WHERE email = sqlc.arg(email) ORDER BY created_at DESC LIMIT 1;

-- name: ListContacts :many
SELECT * FROM contacts
ORDER BY created_at DESC
LIMIT sqlc.arg(page_size) OFFSET sqlc.arg(page_offset);

-- name: CountContacts :one
SELECT COUNT(*) as count FROM contacts;

-- name: SetContactVerificationToken :exec
UPDATE contacts
SET verification_token = sqlc.arg(token), verification_sent_at = datetime('now')
WHERE id = sqlc.arg(id);

-- name: GetContactByVerificationToken :one
SELECT * FROM contacts WHERE verification_token = sqlc.arg(token);

-- name: VerifyContact :exec
UPDATE contacts
SET email_verified = 1, verified_at = datetime('now'), verification_token = NULL
WHERE id = sqlc.arg(id);

-- name: VerifyContactEmail :exec
UPDATE contacts
SET email_verified = 1, verified_at = datetime('now'), verification_token = NULL
WHERE id = sqlc.arg(id);

-- name: UpdateContactStatus :exec
UPDATE contacts
SET status = sqlc.arg(status), updated_at = datetime('now')
WHERE id = sqlc.arg(id);

-- name: CountContactsByStatus :one
SELECT COUNT(*) as count FROM contacts WHERE status = sqlc.arg(status);

-- name: ManuallyVerifyContact :exec
UPDATE contacts
SET email_verified = 1, verified_at = datetime('now'), verification_token = NULL
WHERE id = sqlc.arg(id);

-- name: ResubscribeContact :exec
UPDATE contacts
SET unsubscribed_at = NULL, updated_at = datetime('now')
WHERE id = sqlc.arg(id);

-- name: BlockContact :exec
UPDATE contacts
SET blocked_at = datetime('now'), updated_at = datetime('now')
WHERE id = sqlc.arg(id);

-- name: UnblockContact :exec
UPDATE contacts
SET blocked_at = NULL, updated_at = datetime('now')
WHERE id = sqlc.arg(id);

-- name: ListContactsByOrg :many
SELECT * FROM contacts
WHERE org_id = sqlc.arg(org_id)
ORDER BY created_at DESC
LIMIT sqlc.arg(page_size) OFFSET sqlc.arg(page_offset);

-- name: GetContactByOrgID :one
SELECT * FROM contacts
WHERE id = sqlc.arg(id) AND org_id = sqlc.arg(org_id);

-- name: CountContactsByOrg :one
SELECT COUNT(*) as count FROM contacts WHERE org_id = sqlc.arg(org_id);

-- name: UpdateSDKContact :one
UPDATE contacts
SET name = COALESCE(NULLIF(sqlc.arg(name), ''), name),
    updated_at = datetime('now')
WHERE id = sqlc.arg(id) AND org_id = sqlc.arg(org_id)
RETURNING *;

-- name: GlobalUnsubscribeByOrgAndEmail :exec
UPDATE contacts
SET unsubscribed_at = datetime('now'), status = 'unsubscribed', updated_at = datetime('now')
WHERE org_id = sqlc.arg(org_id) AND email = sqlc.arg(email);

-- name: DeleteContact :exec
DELETE FROM contacts WHERE id = sqlc.arg(id);

-- name: UpdateContactGDPR :exec
UPDATE contacts
SET gdpr_consent = sqlc.arg(consent),
    gdpr_consent_at = CASE WHEN sqlc.arg(consent) = 1 THEN datetime('now') ELSE NULL END,
    updated_at = datetime('now')
WHERE id = sqlc.arg(id);

-- name: UpdateContactName :exec
UPDATE contacts
SET name = sqlc.arg(name), updated_at = datetime('now')
WHERE id = sqlc.arg(id);

-- name: CountUnconfirmedContactsOlderThan :one
SELECT COUNT(*) as count FROM contacts
WHERE org_id = sqlc.arg(org_id)
AND email_verified = 0
AND datetime(created_at) < datetime('now', sqlc.arg(days_modifier));

-- name: DeleteUnconfirmedContactsOlderThan :execrows
DELETE FROM contacts
WHERE org_id = sqlc.arg(org_id)
AND email_verified = 0
AND datetime(created_at) < datetime('now', sqlc.arg(days_modifier));

-- name: CountInactiveContacts90Days :one
SELECT COUNT(*) as count FROM contacts c
WHERE c.org_id = sqlc.arg(org_id)
AND c.email_verified = 1
AND c.unsubscribed_at IS NULL
AND NOT EXISTS (
    SELECT 1 FROM email_queue eq
    WHERE eq.contact_id = c.id
    AND eq.opened_at > datetime('now', '-90 days')
);

-- name: DeleteInactiveContacts90Days :execrows
DELETE FROM contacts
WHERE contacts.org_id = sqlc.arg(org_id)
AND contacts.id IN (
    SELECT c.id FROM contacts c
    WHERE c.org_id = sqlc.arg(org_id)
    AND c.email_verified = 1
    AND c.unsubscribed_at IS NULL
    AND NOT EXISTS (
        SELECT 1 FROM email_queue eq
        WHERE eq.contact_id = c.id
        AND eq.opened_at > datetime('now', '-90 days')
    )
);

-- ============================================
-- DASHBOARD STATS QUERIES
-- ============================================

-- name: GetDashboardSubscriberStats :one
SELECT
    COUNT(*) as total,
    SUM(CASE WHEN datetime(created_at) > datetime('now', '-30 days') THEN 1 ELSE 0 END) as new_30d,
    SUM(CASE WHEN datetime(created_at) > datetime('now', '-7 days') THEN 1 ELSE 0 END) as new_7d
FROM contacts WHERE org_id = sqlc.arg(org_id);

-- name: GetDashboardSubscriberGrowth :one
SELECT
    COALESCE(
        (SELECT COUNT(*) FROM contacts c1 WHERE c1.org_id = sqlc.arg(org_id) AND datetime(c1.created_at) > datetime('now', '-30 days')),
        0
    ) as current_period,
    COALESCE(
        (SELECT COUNT(*) FROM contacts c2 WHERE c2.org_id = sqlc.arg(org_id) AND datetime(c2.created_at) BETWEEN datetime('now', '-60 days') AND datetime('now', '-30 days')),
        0
    ) as previous_period;

-- ============================================
-- GDPR EXPORT QUERIES
-- ============================================

-- name: GetContactCampaignSends :many
SELECT
    cs.id,
    cs.campaign_id,
    ec.name as campaign_name,
    cs.sent_at,
    cs.opened_at,
    cs.clicked_at
FROM campaign_sends cs
LEFT JOIN email_campaigns ec ON ec.id = cs.campaign_id
WHERE cs.contact_id = sqlc.arg(contact_id)
ORDER BY cs.sent_at DESC;

-- name: GetContactEmailClicks :many
SELECT
    ec.id,
    ec.link_url,
    ec.link_name,
    ec.clicked_at,
    ec.user_agent,
    ec.ip_address
FROM email_clicks ec
WHERE ec.contact_id = sqlc.arg(contact_id)
ORDER BY ec.clicked_at DESC;

-- name: GetContactTransactionalSends :many
SELECT
    ts.id,
    ts.template_id,
    te.name as template_name,
    ts.to_email,
    ts.status,
    ts.opened_at,
    ts.clicked_at,
    ts.created_at
FROM transactional_sends ts
LEFT JOIN transactional_emails te ON te.id = ts.template_id
WHERE ts.contact_id = sqlc.arg(contact_id)
ORDER BY ts.created_at DESC;

-- name: GetContactSequenceEmails :many
SELECT
    eq.id,
    et.sequence_id,
    es.name as sequence_name,
    et.position as step_number,
    et.subject,
    eq.status,
    eq.sent_at,
    eq.opened_at,
    eq.clicked_at
FROM email_queue eq
LEFT JOIN email_templates et ON et.id = eq.template_id
LEFT JOIN email_sequences es ON es.id = et.sequence_id
WHERE eq.contact_id = sqlc.arg(contact_id)
ORDER BY eq.sent_at DESC;
