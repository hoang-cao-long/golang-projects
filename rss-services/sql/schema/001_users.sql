-- +goose Up
CREATE TABLE users
(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL
);

-- CREATE TABLE users
-- (
--     id INTEGER PRIMARY KEY,
--     created_at TIMESTAMP NOT NULL,
--     updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
--     name TEXT NOT NULL
-- );

-- +goose Down
DROP TABLE users;
