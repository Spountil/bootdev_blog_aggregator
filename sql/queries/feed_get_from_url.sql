-- name: GetFeedsIDFromURL :one
SELECT id
FROM feed
WHERE url = $1;