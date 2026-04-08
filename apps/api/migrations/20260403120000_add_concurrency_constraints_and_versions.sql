-- +goose Up
-- +goose StatementBegin
ALTER TABLE league_management.organizations
    DROP CONSTRAINT IF EXISTS organizations_name_key;

CREATE UNIQUE INDEX IF NOT EXISTS organizations_user_id_lower_name_idx
    ON league_management.organizations (user_id, lower(name));

ALTER TABLE league_management.leagues
    DROP CONSTRAINT IF EXISTS leagues_name_key;

CREATE UNIQUE INDEX IF NOT EXISTS leagues_organization_id_lower_name_idx
    ON league_management.leagues (organization_id, lower(name));

ALTER TABLE league_management.seasons
    ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;

UPDATE league_management.seasons
SET version = 1
WHERE version IS NULL;

WITH ranked_active_seasons AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY league_id
               ORDER BY updated_at DESC, created_at DESC, id DESC
           ) AS active_rank
    FROM league_management.seasons
    WHERE status IN ('pending', 'planned', 'in_progress', 'paused')
)
UPDATE league_management.seasons s
SET status = 'undefined'
FROM ranked_active_seasons ras
WHERE s.id = ras.id
  AND ras.active_rank > 1;

CREATE UNIQUE INDEX IF NOT EXISTS seasons_one_active_per_league_idx
    ON league_management.seasons (league_id)
    WHERE status IN ('pending', 'planned', 'in_progress', 'paused');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS league_management.seasons_one_active_per_league_idx;

ALTER TABLE league_management.seasons
    DROP COLUMN IF EXISTS version;

DROP INDEX IF EXISTS league_management.leagues_organization_id_lower_name_idx;

ALTER TABLE league_management.leagues
    ADD CONSTRAINT leagues_name_key UNIQUE (name);

DROP INDEX IF EXISTS league_management.organizations_user_id_lower_name_idx;

ALTER TABLE league_management.organizations
    ADD CONSTRAINT organizations_name_key UNIQUE (name);
-- +goose StatementEnd
