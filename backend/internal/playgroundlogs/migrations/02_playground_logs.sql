CREATE TABLE playground_logs (
    id TEXT NOT NULL,
    input TEXT NOT NULL,
    policy TEXT NOT NULL,
    result TEXT NOT NULL,
    coverage TEXT NOT NULL,
    "timestamp" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
