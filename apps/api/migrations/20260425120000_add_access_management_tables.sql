-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.access_role_assignments
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES league_management.users (id),
    role       varchar(100) NOT NULL,
    scope_type varchar(100) NOT NULL,
    scope_id   varchar(255) NOT NULL,
    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_access_role_assignments_unique
    ON league_management.access_role_assignments (user_id, role, scope_type, scope_id);

CREATE INDEX IF NOT EXISTS idx_access_role_assignments_scope
    ON league_management.access_role_assignments (scope_type, scope_id, role);

CREATE TABLE IF NOT EXISTS league_management.access_resource_relationships
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_type varchar(100) NOT NULL,
    resource_id   varchar(255) NOT NULL,
    relation      varchar(100) NOT NULL,
    subject_type  varchar(100) NOT NULL,
    subject_id    varchar(255) NOT NULL,
    created_at    timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_access_resource_relationships_unique
    ON league_management.access_resource_relationships (resource_type, resource_id, relation, subject_type, subject_id);

CREATE INDEX IF NOT EXISTS idx_access_resource_relationships_lookup
    ON league_management.access_resource_relationships (resource_type, resource_id, relation);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS league_management.idx_access_resource_relationships_lookup;
DROP INDEX IF EXISTS league_management.idx_access_resource_relationships_unique;
DROP TABLE IF EXISTS league_management.access_resource_relationships;
DROP INDEX IF EXISTS league_management.idx_access_role_assignments_scope;
DROP INDEX IF EXISTS league_management.idx_access_role_assignments_unique;
DROP TABLE IF EXISTS league_management.access_role_assignments;
-- +goose StatementEnd
