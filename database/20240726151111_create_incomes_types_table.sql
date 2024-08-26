-- +goose Up
-- +goose StatementBegin
CREATE TABLE `income_types` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT UNIQUE,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `fk_user_id` BIGINT NOT NULL,
    FOREIGN KEY (`fk_user_id`) REFERENCES `users`(`id`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `income_types`;
-- +goose StatementEnd
