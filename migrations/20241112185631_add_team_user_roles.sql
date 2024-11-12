-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.team_user_roles
(
    id              UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    team_id         UUID         NOT NULL REFERENCES league_management.teams (id),
    user_id         UUID         NOT NULL REFERENCES league_management.users (id),
    role            varchar(100) NOT NULL,
    created_at      timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX unq_id_team_id_user_id_role ON league_management.team_user_roles (team_id, user_id, role);
CREATE INDEX idx_team_id_user_id_role ON league_management.team_user_roles (team_id, user_id, role);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS league_management.team_user_roles
-- +goose StatementEnd
