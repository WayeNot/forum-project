-- +goose Up
ALTER TABLE sessions ADD COLUMN is_active BOOLEAN DEFAULT 'TRUE';

-- +goose Down
ALTER TABLE sessions DROP COLUMN is_active;