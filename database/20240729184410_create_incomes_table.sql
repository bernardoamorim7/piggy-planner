-- +goose Up
-- +goose StatementBegin
CREATE TABLE `incomes` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT UNIQUE,
    `fk_user_id` BIGINT NOT NULL,
    `amount` DECIMAL(10, 2) NOT NULL,
    `description` TEXT NOT NULL,
    `fk_income_type_id` BIGINT NOT NULL,
    `date` DATE NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`fk_user_id`) REFERENCES `users`(`id`),
    FOREIGN KEY (`fk_income_type_id`) REFERENCES `income_types`(`id`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incomes;
-- +goose StatementEnd