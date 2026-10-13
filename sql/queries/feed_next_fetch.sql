-- name: GetNextFeedFetch :one
SELECT id, url
FROM feed
ORDER BY last_fetched_at DESC NULLS FIRST;