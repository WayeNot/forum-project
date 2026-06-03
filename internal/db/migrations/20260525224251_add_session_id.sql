-- +goose Up
ALTER TABLE sessions ADD COLUMN session_id VARCHAR(55) NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE sessions DROP COLUMN session_id;
