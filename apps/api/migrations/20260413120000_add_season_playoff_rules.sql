-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'season_phase'
          AND n.nspname = 'public'
    ) THEN
        CREATE TYPE season_phase AS ENUM (
            'regular_season',
            'playoffs',
            'completed'
        );
    END IF;
END $$;

ALTER TABLE league_management.seasons
    ADD COLUMN IF NOT EXISTS season_phase season_phase NOT NULL DEFAULT 'regular_season',
    ADD COLUMN IF NOT EXISTS champion_team_id UUID REFERENCES league_management.teams (id);

CREATE TABLE IF NOT EXISTS league_management.season_playoff_rules
(
    season_id                  UUID PRIMARY KEY REFERENCES league_management.seasons (id) ON DELETE CASCADE,
    qualification_type         VARCHAR(64) NOT NULL,
    qualifier_count            INTEGER NOT NULL,
    reseed_each_round          BOOLEAN NOT NULL DEFAULT FALSE,
    third_place_match          BOOLEAN NOT NULL DEFAULT FALSE,
    allow_admin_seed_override  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                 timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                 timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS league_management.season_playoff_round_rules
(
    season_id                        UUID NOT NULL REFERENCES league_management.seasons (id) ON DELETE CASCADE,
    round_order                      INTEGER NOT NULL,
    round_name                       VARCHAR(64) NOT NULL,
    legs                             INTEGER NOT NULL,
    higher_seed_hosts_second_leg     BOOLEAN NOT NULL DEFAULT FALSE,
    tied_aggregate_resolution        VARCHAR(64) NOT NULL,
    created_at                       timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at                       timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (season_id, round_order)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS league_management.season_playoff_round_rules;
DROP TABLE IF EXISTS league_management.season_playoff_rules;

ALTER TABLE league_management.seasons
    DROP COLUMN IF EXISTS champion_team_id,
    DROP COLUMN IF EXISTS season_phase;

DROP TYPE IF EXISTS season_phase;
-- +goose StatementEnd
