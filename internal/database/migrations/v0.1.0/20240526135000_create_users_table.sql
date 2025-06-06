-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id         INTEGER                     NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
    name       TEXT                        NOT NULL,
    email      TEXT                        NOT NULL UNIQUE,
    password   TEXT                        NOT NULL,
    avatar     TEXT                        NOT NULL,
    is_admin   BOOLEAN                     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE          DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd