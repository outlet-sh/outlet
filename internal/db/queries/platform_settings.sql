-- name: GetPlatformSetting :one
SELECT * FROM platform_settings
WHERE key = sqlc.arg(key)
LIMIT 1;

-- name: GetPlatformSettingsByCategory :many
SELECT * FROM platform_settings
WHERE category = sqlc.arg(category)
ORDER BY key;

-- name: ListPlatformSettings :many
SELECT * FROM platform_settings
ORDER BY category, key;

-- name: UpsertPlatformSetting :one
INSERT INTO platform_settings (key, value_encrypted, value_text, description, category, is_sensitive, created_at, updated_at)
VALUES (sqlc.arg(key), sqlc.arg(value_encrypted), sqlc.arg(value_text), sqlc.arg(description), sqlc.arg(category), sqlc.arg(is_sensitive), datetime('now'), datetime('now'))
ON CONFLICT (key) DO UPDATE SET
    value_encrypted = CASE WHEN excluded.value_encrypted IS NOT NULL THEN excluded.value_encrypted ELSE platform_settings.value_encrypted END,
    value_text = CASE WHEN excluded.value_text IS NOT NULL THEN excluded.value_text ELSE platform_settings.value_text END,
    description = COALESCE(excluded.description, platform_settings.description),
    updated_at = datetime('now')
RETURNING *;

-- name: DeletePlatformSetting :exec
DELETE FROM platform_settings
WHERE key = sqlc.arg(key);

-- name: GetPlatformSettingValue :one
SELECT key, value_encrypted, value_text, is_sensitive
FROM platform_settings
WHERE key = sqlc.arg(key)
LIMIT 1;

-- name: GetPlatformSettingsValues :many
SELECT key, value_encrypted, value_text, is_sensitive
FROM platform_settings
WHERE category = sqlc.arg(category)
ORDER BY key;
