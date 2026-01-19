-- name: CreateCustomField :one
INSERT INTO custom_fields (id, list_id, name, field_key, field_type, options, required, default_value, placeholder, sort_order, created_at, updated_at)
VALUES (sqlc.arg(id), sqlc.arg(list_id), sqlc.arg(name), sqlc.arg(field_key), sqlc.arg(field_type), sqlc.arg(options), sqlc.arg(required), sqlc.arg(default_value), sqlc.arg(placeholder), sqlc.arg(sort_order), datetime('now'), datetime('now'))
RETURNING *;

-- name: GetCustomField :one
SELECT * FROM custom_fields WHERE id = sqlc.arg(id);

-- name: GetCustomFieldByKey :one
SELECT * FROM custom_fields WHERE list_id = sqlc.arg(list_id) AND field_key = sqlc.arg(field_key);

-- name: ListCustomFieldsByList :many
SELECT * FROM custom_fields WHERE list_id = sqlc.arg(list_id) ORDER BY sort_order ASC, created_at ASC;

-- name: UpdateCustomField :one
UPDATE custom_fields
SET name = sqlc.arg(name),
    field_key = sqlc.arg(field_key),
    field_type = sqlc.arg(field_type),
    options = sqlc.arg(options),
    required = sqlc.arg(required),
    default_value = sqlc.arg(default_value),
    placeholder = sqlc.arg(placeholder),
    sort_order = sqlc.arg(sort_order),
    updated_at = datetime('now')
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteCustomField :exec
DELETE FROM custom_fields WHERE id = sqlc.arg(id);

-- name: DeleteCustomFieldsByList :exec
DELETE FROM custom_fields WHERE list_id = sqlc.arg(list_id);

-- name: CountCustomFieldsByList :one
SELECT COUNT(*) FROM custom_fields WHERE list_id = sqlc.arg(list_id);

-- Custom Field Values

-- name: CreateCustomFieldValue :one
INSERT INTO custom_field_values (id, subscriber_id, field_id, value, created_at, updated_at)
VALUES (sqlc.arg(id), sqlc.arg(subscriber_id), sqlc.arg(field_id), sqlc.arg(value), datetime('now'), datetime('now'))
RETURNING *;

-- name: UpsertCustomFieldValue :one
INSERT INTO custom_field_values (id, subscriber_id, field_id, value, created_at, updated_at)
VALUES (sqlc.arg(id), sqlc.arg(subscriber_id), sqlc.arg(field_id), sqlc.arg(value), datetime('now'), datetime('now'))
ON CONFLICT(subscriber_id, field_id) DO UPDATE SET
    value = excluded.value,
    updated_at = datetime('now')
RETURNING *;

-- name: GetCustomFieldValue :one
SELECT * FROM custom_field_values WHERE subscriber_id = sqlc.arg(subscriber_id) AND field_id = sqlc.arg(field_id);

-- name: ListCustomFieldValuesBySubscriber :many
SELECT cfv.*, cf.name, cf.field_key, cf.field_type
FROM custom_field_values cfv
JOIN custom_fields cf ON cf.id = cfv.field_id
WHERE cfv.subscriber_id = sqlc.arg(subscriber_id)
ORDER BY cf.sort_order ASC;

-- name: GetSubscriberCustomFieldsForMerge :many
-- Returns field_key -> value pairs for use in email merge tags
SELECT cf.field_key, cfv.value
FROM custom_field_values cfv
JOIN custom_fields cf ON cf.id = cfv.field_id
WHERE cfv.subscriber_id = sqlc.arg(subscriber_id);

-- name: DeleteCustomFieldValue :exec
DELETE FROM custom_field_values WHERE subscriber_id = sqlc.arg(subscriber_id) AND field_id = sqlc.arg(field_id);

-- name: DeleteCustomFieldValuesBySubscriber :exec
DELETE FROM custom_field_values WHERE subscriber_id = sqlc.arg(subscriber_id);

-- name: BulkCreateCustomFieldValues :exec
-- Used when subscribing with multiple field values at once
INSERT INTO custom_field_values (id, subscriber_id, field_id, value, created_at, updated_at)
SELECT
    lower(hex(randomblob(16))),
    sqlc.arg(subscriber_id),
    cf.id,
    json_extract(sqlc.arg(values_json), '$.' || cf.field_key),
    datetime('now'),
    datetime('now')
FROM custom_fields cf
WHERE cf.list_id = sqlc.arg(list_id)
  AND json_extract(sqlc.arg(values_json), '$.' || cf.field_key) IS NOT NULL
ON CONFLICT(subscriber_id, field_id) DO UPDATE SET
    value = excluded.value,
    updated_at = datetime('now');
