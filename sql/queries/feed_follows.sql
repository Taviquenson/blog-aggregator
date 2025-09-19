-- name: CreateFeedFollow :one
INSERT INTO feed_follows (
    id, created_at, updated_at, user_id, feed_id
)
values($1,$2,$3,$4,$5)
RETURNING *, (
    SELECT name
    FROM users
    WHERE users.id = feed_follows.user_id
) AS user_name, (
    SELECT name
    FROM feeds
    WHERE feeds.id = feed_follows.feed_id
) AS feed_name;

-- name: GetFeedFollowsForUser :many
SELECT *, (
    SELECT name
    FROM users
    WHERE users.id = feed_follows.user_id
) AS user_name, (
    SELECT name
    FROM feeds
    WHERE feeds.id = feed_follows.feed_id
) AS feed_name
FROM feed_follows
WHERE feed_follows.user_id = (
    SELECT id
    FROM users
    WHERE users.name = $1
);

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE feed_id = $1 AND user_id = $2;
