-- +goose Up
ALTER TABLE experiences ADD COLUMN IF NOT EXISTS logo TEXT;

-- +goose Down
ALTER TABLE experiences DROP COLUMN IF EXISTS logo;
