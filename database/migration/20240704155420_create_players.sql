-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS players (
   id VARCHAR(50) PRIMARY KEY,
   player_name VARCHAR(50) NOT NULL,
   position VARCHAR(50) NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players;
-- +goose StatementEnd

