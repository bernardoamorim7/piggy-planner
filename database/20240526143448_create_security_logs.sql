-- +goose Up
-- +goose StatementBegin
CREATE TABLE `security_logs` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `fk_user_id` BIGINT NOT NULL,
    `action` ENUM(
        'register',
        'login',
        'logout',
        'password_change',
        'password_reset'
    ) NOT NULL,
    `ip_address` VARCHAR(255) NOT NULL,
    `user_agent` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (fk_user_id) REFERENCES user(id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE `security_logs`;

-- +goose StatementEnd