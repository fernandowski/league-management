package pg

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"league-management/internal/access_management/domain"
	"league-management/internal/shared/app_errors"
	"league-management/internal/shared/database"
)

type AccessRepository struct{}

type accessQuerier interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func NewAccessRepository() *AccessRepository {
	return &AccessRepository{}
}

func (ar *AccessRepository) GrantRole(assignment *domain.RoleAssignment) error {
	return ar.GrantRoleTx(context.Background(), database.GetConnection(), assignment)
}

func (ar *AccessRepository) GrantRoleTx(ctx context.Context, tx accessQuerier, assignment *domain.RoleAssignment) error {
	sql := `
		INSERT INTO league_management.access_role_assignments (user_id, role, scope_type, scope_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, role, scope_type, scope_id) DO NOTHING`

	_, err := tx.Exec(ctx, sql, assignment.UserID, string(assignment.Role), string(assignment.Scope.Type), assignment.Scope.ID)
	return err
}

func (ar *AccessRepository) RevokeRole(assignment *domain.RoleAssignment) error {
	sql := `
		DELETE FROM league_management.access_role_assignments
		WHERE user_id = $1 AND role = $2 AND scope_type = $3 AND scope_id = $4`

	_, err := database.GetConnection().Exec(
		context.Background(),
		sql,
		assignment.UserID,
		string(assignment.Role),
		string(assignment.Scope.Type),
		assignment.Scope.ID,
	)

	return err
}

func (ar *AccessRepository) HasRole(userID string, role domain.Role, scope domain.Scope) (bool, error) {
	sql := `
		SELECT EXISTS (
			SELECT 1
			FROM league_management.access_role_assignments
			WHERE user_id = $1 AND role = $2 AND scope_type = $3 AND scope_id = $4
		)`

	var exists bool
	err := database.GetConnection().QueryRow(context.Background(), sql, userID, string(role), string(scope.Type), scope.ID).Scan(&exists)
	return exists, err
}

func (ar *AccessRepository) ListGrants(userID, scopeType, scopeID string) ([]domain.RoleAssignment, error) {
	sql := `
		SELECT id, user_id, role, scope_type, scope_id
		FROM league_management.access_role_assignments
		WHERE ($1 = '' OR user_id = $1)
		  AND ($2 = '' OR scope_type = $2)
		  AND ($3 = '' OR scope_id = $3)
		ORDER BY created_at DESC`

	rows, err := database.GetConnection().Query(context.Background(), sql, userID, scopeType, scopeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grants := []domain.RoleAssignment{}
	for rows.Next() {
		var grant domain.RoleAssignment
		var role string
		var scopeTypeValue string
		var scopeID string

		if err := rows.Scan(&grant.ID, &grant.UserID, &role, &scopeTypeValue, &scopeID); err != nil {
			return nil, err
		}

		parsedRole, err := domain.ParseRole(role)
		if err != nil {
			return nil, err
		}

		parsedScopeType, err := domain.ParseScopeType(scopeTypeValue)
		if err != nil {
			return nil, err
		}

		scope, err := domain.NewScope(parsedScopeType, scopeID)
		if err != nil {
			return nil, err
		}

		grant.Role = parsedRole
		grant.Scope = scope
		grants = append(grants, grant)
	}

	return grants, rows.Err()
}

func (ar *AccessRepository) RecordResourceRelationship(relationship *domain.ResourceRelationship) error {
	return ar.RecordResourceRelationshipTx(context.Background(), database.GetConnection(), relationship)
}

func (ar *AccessRepository) RecordResourceRelationshipTx(ctx context.Context, tx accessQuerier, relationship *domain.ResourceRelationship) error {
	sql := `
		INSERT INTO league_management.access_resource_relationships
			(resource_type, resource_id, relation, subject_type, subject_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (resource_type, resource_id, relation, subject_type, subject_id) DO NOTHING`

	_, err := tx.Exec(
		ctx,
		sql,
		string(relationship.ResourceType),
		relationship.ResourceID,
		string(relationship.Relation),
		string(relationship.SubjectType),
		relationship.SubjectID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return app_errors.ErrDuplicateResource
		}
	}

	return err
}

func (ar *AccessRepository) FindParent(resourceType domain.ResourceType, resourceID string) (*domain.ResourceRelationship, error) {
	sql := `
		SELECT resource_type, resource_id, relation, subject_type, subject_id
		FROM league_management.access_resource_relationships
		WHERE resource_type = $1 AND resource_id = $2 AND relation = $3
		LIMIT 1`

	var resourceTypeValue, relationValue, subjectTypeValue string
	var relationship domain.ResourceRelationship

	err := database.GetConnection().QueryRow(context.Background(), sql, string(resourceType), resourceID, string(domain.RelationParent)).
		Scan(&resourceTypeValue, &relationship.ResourceID, &relationValue, &subjectTypeValue, &relationship.SubjectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	relationship.ResourceType = domain.ResourceType(resourceTypeValue)
	relationship.Relation = domain.Relation(relationValue)
	relationship.SubjectType = domain.SubjectType(subjectTypeValue)

	return &relationship, nil
}

func (ar *AccessRepository) HasOwnerRelationship(organizationID, userID string) (bool, error) {
	sql := `
		SELECT EXISTS (
			SELECT 1
			FROM league_management.access_resource_relationships
			WHERE resource_type = $1
			  AND resource_id = $2
			  AND relation = $3
			  AND subject_type = $4
			  AND subject_id = $5
		)`

	var exists bool
	err := database.GetConnection().QueryRow(
		context.Background(),
		sql,
		string(domain.ResourceOrganization),
		organizationID,
		string(domain.RelationOwner),
		string(domain.SubjectUser),
		userID,
	).Scan(&exists)

	return exists, err
}
