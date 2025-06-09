-- +goose Up
ALTER TABLE users
ADD CONSTRAINT name_not_blank
CHECK (BTRIM(name) <> '');

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT name_not_blank;
