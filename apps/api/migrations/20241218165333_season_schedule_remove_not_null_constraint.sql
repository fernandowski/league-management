-- +goose Up
-- +goose StatementBegin
ALTER TABLE league_management.season_schedules ALTER COLUMN home_team_id DROP NOT NULL;
ALTER TABLE league_management.season_schedules ALTER COLUMN away_team_id DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE league_management.season_schedules ALTER COLUMN home_team_id SET NOT NULL;
ALTER TABLE league_management.season_schedules ALTER COLUMN away_team_id SET NOT NULL;
-- +goose StatementEnd
