-- name: MarkFeedFetched :exec
UPDATE feed
SET updated_at = $1, last_fetched_at = $2
WHERE id = $3;