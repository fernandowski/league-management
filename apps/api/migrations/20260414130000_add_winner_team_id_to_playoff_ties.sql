-- +goose Up
-- +goose StatementBegin
ALTER TABLE league_management.season_playoff_ties
    ADD COLUMN IF NOT EXISTS winner_team_id UUID REFERENCES league_management.teams (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE league_management.season_playoff_ties
    DROP COLUMN IF EXISTS winner_team_id;
-- +goose StatementEnd
