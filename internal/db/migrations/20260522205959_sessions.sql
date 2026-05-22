-- +goose Up
CREATE TABLE sessions (
    id INTEGER PRIMARY KEY,
    user_id VARCHAR(60) NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE sessions;