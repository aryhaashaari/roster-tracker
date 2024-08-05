-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    email VARCHAR(50) PRIMARY KEY,
    password VARCHAR(50) NOT NULL,
    code VARCHAR(25) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
