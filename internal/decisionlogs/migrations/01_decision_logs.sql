CREATE TABLE decision_logs (
    decision_id TEXT NOT NULL,
    "path" TEXT NOT NULL,
    input TEXT NOT NULL,
    revision_id TEXT,
    result TEXT NOT NULL,
    "timestamp" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
