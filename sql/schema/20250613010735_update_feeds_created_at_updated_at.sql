-- +goose Up
ALTER TABLE feeds RENAME COLUMN createdAt TO created_at;
ALTER TABLE feeds RENAME COLUMN updatedAt TO updated_at;

-- +goose Down
ALTER TABLE feeds RENAME COLUMN created_at TO createdAt;
ALTER TABLE feeds RENAME COLUMN updated_at TO updatedAt;