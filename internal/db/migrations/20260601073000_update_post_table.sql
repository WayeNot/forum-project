-- +goose Up
ALTER TABLE posts ADD COLUMN tags TEXT;

-- +goose Down
ALTER TABLE posts DROP COLUMN tags;