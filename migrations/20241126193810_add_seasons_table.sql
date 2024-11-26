-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.seasons
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    league_id  UUID         NOT NULL REFERENCES league_management.leagues (id),
    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_seasons_name ON league_management.seasons (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS league_management.seasons
-- +goose StatementEnd
