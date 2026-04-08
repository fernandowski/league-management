-- +goose Up
-- +goose StatementBegin
ALTER TABLE league_management.league_teams
DROP CONSTRAINT league_teams_team_id_fkey;

ALTER TABLE league_management.league_teams
ADD CONSTRAINT league_teams_team_id_fkey
FOREIGN KEY (team_id) REFERENCES league_management.teams (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
