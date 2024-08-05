-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stats (
   id VARCHAR(50) NOT NULL,
   stats_id VARCHAR(50) NOT NULL,
   points VARCHAR(50) NOT NULL,
   assists VARCHAR(50) NOT NULL,
   rebounds VARCHAR(50) NOT NULL,
   fieldGoalPct VARCHAR(50) NOT NULL,
   threePointPct VARCHAR(50) NOT NULL,
   steals VARCHAR(50) NOT NULL,
   blocks VARCHAR(50) NOT NULL,
   turnovers VARCHAR(50) NOT NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stats;
-- +goose StatementEnd