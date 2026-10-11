-- name: GetFeedFollowsForUser :many
WITH user_id AS (
    SELECT id
    FROM users
    WHERE users.name = $1
)
SELECT feed.name
FROM feed_follows
JOIN user_id
    ON feed_follows.user_id = user_id.id
JOIN feed
    ON feed_follows.feed_id = feed.id;