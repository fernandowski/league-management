-- +goose Up
-- +goose StatementBegin
INSERT INTO league_management.access_resource_relationships
    (resource_type, resource_id, relation, subject_type, subject_id)
SELECT
    'organization',
    organizations.id::text,
    'owner',
    'user',
    organizations.user_id::text
FROM league_management.organizations
WHERE organizations.user_id IS NOT NULL
ON CONFLICT (resource_type, resource_id, relation, subject_type, subject_id) DO NOTHING;

INSERT INTO league_management.access_resource_relationships
    (resource_type, resource_id, relation, subject_type, subject_id)
SELECT
    'team',
    organization_teams.team_id::text,
    'parent',
    'organization',
    organization_teams.organization_id::text
FROM league_management.organization_teams
ON CONFLICT (resource_type, resource_id, relation, subject_type, subject_id) DO NOTHING;

INSERT INTO league_management.access_resource_relationships
    (resource_type, resource_id, relation, subject_type, subject_id)
SELECT
    'league',
    leagues.id::text,
    'parent',
    'organization',
    leagues.organization_id::text
FROM league_management.leagues
WHERE leagues.organization_id IS NOT NULL
ON CONFLICT (resource_type, resource_id, relation, subject_type, subject_id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- This is a data backfill. Existing and newly-created facts are indistinguishable,
-- so rolling back must not delete live authorization relationships.
-- +goose StatementEnd
