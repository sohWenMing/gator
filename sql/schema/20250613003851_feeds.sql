-- +goose Up
CREATE TABLE feeds(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL CONSTRAINT unique_url UNIQUE,
    user_id UUID CONSTRAINT user_id_fk REFERENCES users(id)
    ON DELETE CASCADE,
    createdAt TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE feeds;