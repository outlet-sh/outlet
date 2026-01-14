-- name: GetEmailListByOrgAndSlug :one
SELECT * FROM email_lists
WHERE org_id = sqlc.arg(org_id) AND slug = sqlc.arg(slug)
LIMIT 1;

-- name: CreateEmailList :one
INSERT INTO email_lists (
    org_id, name, slug, description, double_optin, created_at, updated_at
) VALUES (
    sqlc.arg(org_id), sqlc.arg(name), sqlc.arg(slug), sqlc.arg(description),
    sqlc.arg(double_optin), datetime('now'), datetime('now')
) RETURNING *;

-- name: ListEmailLists :many
SELECT * FROM email_lists
WHERE org_id = sqlc.arg(org_id)
ORDER BY created_at DESC;

-- name: GetEmailList :one
SELECT * FROM email_lists
WHERE id = sqlc.arg(id);

-- name: UpdateEmailList :one
UPDATE email_lists
SET name = COALESCE(sqlc.arg(name), name),
    description = COALESCE(sqlc.arg(description), description),
    double_optin = COALESCE(sqlc.arg(double_optin), double_optin),
    confirmation_email_subject = COALESCE(NULLIF(sqlc.arg(confirmation_subject), ''), confirmation_email_subject),
    confirmation_email_body = COALESCE(NULLIF(sqlc.arg(confirmation_body), ''), confirmation_email_body),
    updated_at = datetime('now')
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteEmailList :exec
DELETE FROM email_lists
WHERE id = sqlc.arg(id);

-- name: SubscribeToList :one
INSERT INTO list_subscribers (
    id, list_id, contact_id, status, subscribed_at
) VALUES (
    sqlc.arg(id), sqlc.arg(list_id), sqlc.arg(contact_id), 'active', datetime('now')
)
ON CONFLICT (list_id, contact_id) DO UPDATE
SET status = 'active', unsubscribed_at = NULL
RETURNING *;

-- name: UnsubscribeFromList :exec
UPDATE list_subscribers
SET status = 'unsubscribed', unsubscribed_at = datetime('now')
WHERE list_id = sqlc.arg(list_id) AND contact_id = sqlc.arg(contact_id);

-- name: GetListSubscriber :one
SELECT * FROM list_subscribers
WHERE list_id = sqlc.arg(list_id) AND contact_id = sqlc.arg(contact_id);

-- name: ListListSubscribers :many
SELECT ls.*, c.email, c.name
FROM list_subscribers ls
JOIN contacts c ON c.id = ls.contact_id
WHERE ls.list_id = sqlc.arg(list_id)
  AND (sqlc.arg(filter_status) IS NULL OR sqlc.arg(filter_status) = '' OR ls.status = sqlc.arg(filter_status))
ORDER BY ls.subscribed_at DESC
LIMIT sqlc.arg(page_size) OFFSET sqlc.arg(page_offset);

-- name: CountListSubscribers :one
SELECT COUNT(*) as count FROM list_subscribers
WHERE list_id = sqlc.arg(list_id)
  AND (sqlc.arg(filter_status) IS NULL OR sqlc.arg(filter_status) = '' OR status = sqlc.arg(filter_status));

-- name: SubscribeToListPending :one
INSERT INTO list_subscribers (
    id, list_id, contact_id, status, verification_token, verification_sent_at
) VALUES (
    sqlc.arg(id), sqlc.arg(list_id), sqlc.arg(contact_id), 'pending',
    sqlc.arg(verification_token), datetime('now')
)
ON CONFLICT (list_id, contact_id) DO UPDATE
SET status = CASE
    WHEN list_subscribers.status = 'active' THEN 'active'
    ELSE 'pending'
END,
verification_token = CASE
    WHEN list_subscribers.status = 'active' THEN list_subscribers.verification_token
    ELSE EXCLUDED.verification_token
END,
verification_sent_at = CASE
    WHEN list_subscribers.status = 'active' THEN list_subscribers.verification_sent_at
    ELSE datetime('now')
END,
unsubscribed_at = NULL
RETURNING *;

-- name: ConfirmListSubscription :one
UPDATE list_subscribers
SET status = 'active', verified_at = datetime('now'), verification_token = NULL
WHERE verification_token = sqlc.arg(token) AND status = 'pending'
RETURNING *;

-- name: GetListSubscriberByToken :one
SELECT ls.*, el.name as list_name, el.org_id
FROM list_subscribers ls
JOIN email_lists el ON el.id = ls.list_id
WHERE ls.verification_token = sqlc.arg(token);

-- name: RemoveSubscriberFromList :exec
DELETE FROM list_subscribers
WHERE list_id = sqlc.arg(list_id) AND contact_id = sqlc.arg(contact_id);

-- name: GetListSubscriberByID :one
SELECT ls.*, c.email, c.name
FROM list_subscribers ls
JOIN contacts c ON c.id = ls.contact_id
WHERE ls.id = sqlc.arg(id);

-- name: CountActiveSubscribers :one
SELECT COUNT(*) as count FROM list_subscribers
WHERE list_id = sqlc.arg(list_id) AND status = 'active';

-- name: CheckListSubscription :one
SELECT COUNT(*) as count FROM list_subscribers
WHERE list_id = sqlc.arg(list_id) AND contact_id = sqlc.arg(contact_id);

-- name: CreateListSubscriber :one
INSERT INTO list_subscribers (list_id, contact_id, status, subscribed_at)
VALUES (sqlc.arg(list_id), sqlc.arg(contact_id), sqlc.arg(status), datetime('now'))
ON CONFLICT (list_id, contact_id) DO NOTHING
RETURNING *;
