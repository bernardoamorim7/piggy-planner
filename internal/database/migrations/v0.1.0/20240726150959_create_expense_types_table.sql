-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS expense_types (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    fk_user_id INTEGER NOT NULL,
    FOREIGN KEY (fk_user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS expense_types;
-- +goose StatementEnd