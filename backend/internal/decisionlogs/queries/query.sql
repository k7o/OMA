-- name: GetDecisionLog :one
SELECT * FROM decision_logs
WHERE decision_id = ? LIMIT 1;

-- name: ListDecisionLogs :many
SELECT * FROM decision_logs
ORDER BY "timestamp" DESC;

-- name: CreateDecisionLog :one
INSERT INTO decision_logs (
  decision_id, path, input, result, timestamp
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;