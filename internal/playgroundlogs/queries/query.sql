-- name: GetPlaygroundLog :one
SELECT * FROM playground_logs
WHERE id = ? LIMIT 1;

-- name: ListPlaygroundlogs :many
SELECT * FROM playground_logs
ORDER BY "timestamp" DESC;

-- name: CreatePlaygroundLog :one
INSERT INTO playground_logs (
  id, input, policy, result, coverage, timestamp
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;