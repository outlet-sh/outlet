-- name: AddContactTag :one
INSERT INTO contact_tags (contact_id, tag, created_at)
VALUES (sqlc.arg(contact_id), sqlc.arg(tag), datetime('now'))
ON CONFLICT (contact_id, tag) DO NOTHING
RETURNING *;

-- name: RemoveContactTag :exec
DELETE FROM contact_tags
WHERE contact_id = sqlc.arg(contact_id) AND tag = sqlc.arg(tag);

-- name: GetContactTags :many
SELECT id, contact_id, tag, created_at
FROM contact_tags
WHERE contact_id = sqlc.arg(contact_id)
ORDER BY created_at;

-- name: HasContactTag :one
SELECT CASE WHEN COUNT(*) > 0 THEN 1 ELSE 0 END as has_tag
FROM contact_tags WHERE contact_id = sqlc.arg(contact_id) AND tag = sqlc.arg(tag);

-- name: GetContactsByTag :many
SELECT c.*
FROM contacts c
JOIN contact_tags ct ON ct.contact_id = c.id
WHERE ct.tag = sqlc.arg(tag)
ORDER BY c.created_at DESC
LIMIT sqlc.arg(page_size) OFFSET sqlc.arg(page_offset);

-- name: CountContactsByTag :one
SELECT COUNT(DISTINCT contact_id) as count
FROM contact_tags
WHERE tag = sqlc.arg(tag);

-- name: ListAllContactTags :many
SELECT tag, COUNT(*) as count
FROM contact_tags
GROUP BY tag
ORDER BY count DESC;
