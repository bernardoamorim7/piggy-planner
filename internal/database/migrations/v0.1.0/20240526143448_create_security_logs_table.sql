-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS security_logs (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
    fk_user_id INTEGER NOT NULL,
    action TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (fk_user_id) REFERENCES users(id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS security_logs;
-- +goose StatementEnd