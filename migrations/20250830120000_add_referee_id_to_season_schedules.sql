-- +goose Up
-- +goose StatementBegin
ALTER TABLE league_management.season_schedules
ADD COLUMN IF NOT EXISTS referee_id VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE league_management.season_schedules
DROP COLUMN IF EXISTS referee_id;
-- +goose StatementEnd
