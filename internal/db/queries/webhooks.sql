-- name: CreateWebhook :one
INSERT INTO webhooks (id, org_id, url, secret, events, active, created_at, updated_at)
VALUES (sqlc.arg(id), sqlc.arg(org_id), sqlc.arg(url), sqlc.arg(secret), sqlc.arg(events), sqlc.arg(active), datetime('now'), datetime('now'))
RETURNING *;

-- name: GetWebhook :one
SELECT * FROM webhooks
WHERE id = sqlc.arg(id) AND org_id = sqlc.arg(org_id);

-- name: GetWebhookByID :one
SELECT * FROM webhooks
WHERE id = sqlc.arg(id);

-- name: ListWebhooks :many
SELECT * FROM webhooks
WHERE org_id = sqlc.arg(org_id)
ORDER BY created_at DESC;

-- name: UpdateWebhook :one
UPDATE webhooks
SET url = COALESCE(NULLIF(sqlc.arg(url), ''), url),
    events = COALESCE(NULLIF(sqlc.arg(events), ''), events),
    active = sqlc.arg(active),
    updated_at = datetime('now')
WHERE id = sqlc.arg(id) AND org_id = sqlc.arg(org_id)
RETURNING *;

-- name: DeleteWebhook :exec
DELETE FROM webhooks
WHERE id = sqlc.arg(id) AND org_id = sqlc.arg(org_id);

-- name: UpdateWebhookDeliveryStats :exec
UPDATE webhooks
SET deliveries_total = deliveries_total + 1,
    deliveries_success = CASE WHEN sqlc.arg(success) = 1 THEN deliveries_success + 1 ELSE deliveries_success END,
    deliveries_failed = CASE WHEN sqlc.arg(success) = 0 THEN deliveries_failed + 1 ELSE deliveries_failed END,
    last_delivery_at = datetime('now'),
    last_status = sqlc.arg(last_status)
WHERE id = sqlc.arg(id);

-- name: CreateWebhookLog :one
INSERT INTO webhook_logs (id, webhook_id, event, payload, status_code, response, error, duration_ms, delivered_at)
VALUES (sqlc.arg(id), sqlc.arg(webhook_id), sqlc.arg(event), sqlc.arg(payload), sqlc.arg(status_code), sqlc.arg(response), sqlc.arg(error), sqlc.arg(duration_ms), datetime('now'))
RETURNING *;

-- name: ListWebhookLogs :many
SELECT * FROM webhook_logs
WHERE webhook_id = sqlc.arg(webhook_id)
ORDER BY delivered_at DESC
LIMIT sqlc.arg(limit_count);

-- name: CountWebhookLogs :one
SELECT COUNT(*) as count FROM webhook_logs
WHERE webhook_id = sqlc.arg(webhook_id);
