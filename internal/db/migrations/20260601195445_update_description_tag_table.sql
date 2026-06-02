-- +goose Up
ALTER TABLE tags ADD COLUMN description TEXT;

-- +goose Down
ALTER TABLE tags DROP COLUMN description;