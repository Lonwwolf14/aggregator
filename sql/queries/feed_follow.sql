-- name: CreateFeedFollow :one
WITH new_follow AS (
    INSERT INTO feed_follow (id, user_id, feed_id)
    VALUES ($1, $2, $3)
    RETURNING *
)
SELECT 
    new_follow.*,
    users.name as user_name,
    feeds.name as feed_name
FROM new_follow
JOIN users ON users.id = new_follow.user_id
JOIN feeds ON feeds.id = new_follow.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follow WHERE user_id = $1;