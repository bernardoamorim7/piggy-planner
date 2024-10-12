-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS incomes
(
    id                INTEGER                     NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
    fk_user_id        INTEGER                     NOT NULL,
    amount            REAL                        NOT NULL,
    description       TEXT                        NOT NULL,
    fk_income_type_id INTEGER                     NOT NULL,
    date              DATE                        NOT NULL,
    created_at        TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITHOUT TIME ZONE          DEFAULT NULL,
    FOREIGN KEY (fk_user_id) REFERENCES users (id),
    FOREIGN KEY (fk_income_type_id) REFERENCES income_types (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS incomes;
-- +goose StatementEnd