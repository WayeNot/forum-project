-- +goose Up
ALTER TABLE sessions ADD COLUMN session_id VARCHAR(55) NOT NULL;

-- +goose Down
ALTER TABLE sessions DROP COLUMN session_id;