-- +goose Up
ALTER TABLE posts DROP COLUMN tags;
ALTER TABLE posts ADD COLUMN tags TEXT[] DEFAULT '{}';

-- +goose Down
ALTER TABLE posts ALTER COLUMN tags TYPE TEXT USING array_to_string(tags, ',') DEFAULT '';