-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    CONSTRAINT UQ_name UNIQUE (name)
);

-- +goose Down
DROP TABLE users;