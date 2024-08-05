-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS physique (
   id VARCHAR(50) PRIMARY KEY,
   height VARCHAR(50) NOT NULL,
   weight VARCHAR(50) NOT NULL,
   age VARCHAR(50) NOT NULL,
   wingspan VARCHAR(50) NOT NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS physique;
-- +goose StatementEnd