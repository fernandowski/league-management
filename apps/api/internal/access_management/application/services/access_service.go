package services

import (
	"errors"
	"league-management/internal/access_management/domain"
)

type AccessService struct {
	repository accessRepository
}

type accessRepository interface {
	GrantRole(*domain.RoleAssignment) error
	RevokeRole(*domain.RoleAssignment) error
	HasRole(userID string, role domain.Role, scope domain.Scope) (bool, error)
	ListGrants(userID, scopeType, scopeID string) ([]domain.RoleAssignment, error)
	RecordResourceRelationship(*domain.ResourceRelationship) error
	FindParent(resourceType domain.ResourceType, resourceID string) (*domain.ResourceRelationship, error)
	HasOwnerRelationship(organizationID, userID string) (bool, error)
}

func NewAccessService(repository accessRepository) *AccessService {
	return &AccessService{repository: repository}
}

// BootstrapSuperAdmin is intentionally reserved for trusted local administration.
// HTTP callers must use the normal authorization flow.
func (as *AccessService) BootstrapSuperAdmin(userID string) error {
	scope, err := domain.NewScope(domain.ScopePlatform, domain.PlatformScopeID)
	if err != nil {
		return err
	}

	assignment, err := domain.NewRoleAssignment(userID, domain.RoleSuperAdmin, scope)
	if err != nil {
		return err
	}

	return as.repository.GrantRole(assignment)
}

func (as *AccessService) GrantRole(actorUserID, targetUserID, roleValue, scopeTypeValue, scopeID string) error {
	assignment, err := newAssignment(targetUserID, roleValue, scopeTypeValue, scopeID)
	if err != nil {
		return err
	}

	if err := as.authorizeGrant(actorUserID, assignment.Role, assignment.Scope); err != nil {
		return err
	}

	return as.repository.GrantRole(assignment)
}

func (as *AccessService) RevokeRole(actorUserID, targetUserID, roleValue, scopeTypeValue, scopeID string) error {
	assignment, err := newAssignment(targetUserID, roleValue, scopeTypeValue, scopeID)
	if err != nil {
		return err
	}

	if err := as.authorizeGrant(actorUserID, assignment.Role, assignment.Scope); err != nil {
		return err
	}

	return as.repository.RevokeRole(assignment)
}

func (as *AccessService) HasRole(userID, roleValue, scopeTypeValue, scopeID string) (bool, error) {
	return as.hasRole(userID, roleValue, scopeTypeValue, scopeID)
}

func (as *AccessService) CheckRole(actorUserID, userID, roleValue, scopeTypeValue, scopeID string) (bool, error) {
	if actorUserID == "" {
		return false, errors.New("actor user id cannot be empty")
	}
	if actorUserID != userID {
		if err := as.authorizeScopeRead(actorUserID, scopeTypeValue, scopeID); err != nil {
			return false, err
		}
	}

	return as.hasRole(userID, roleValue, scopeTypeValue, scopeID)
}

func (as *AccessService) hasRole(userID, roleValue, scopeTypeValue, scopeID string) (bool, error) {
	role, err := domain.ParseRole(roleValue)
	if err != nil {
		return false, err
	}

	scopeType, err := domain.ParseScopeType(scopeTypeValue)
	if err != nil {
		return false, err
	}

	scope, err := domain.NewScope(scopeType, scopeID)
	if err != nil {
		return false, err
	}

	return as.repository.HasRole(userID, role, scope)
}

func (as *AccessService) ListGrants(userID, scopeType, scopeID string) ([]domain.RoleAssignment, error) {
	return as.repository.ListGrants(userID, scopeType, scopeID)
}

func (as *AccessService) ListGrantsForActor(actorUserID, userID, scopeType, scopeID string) ([]domain.RoleAssignment, error) {
	if actorUserID == "" {
		return nil, errors.New("actor user id cannot be empty")
	}
	if actorUserID == userID && scopeType == "" && scopeID == "" {
		return as.repository.ListGrants(userID, scopeType, scopeID)
	}

	if scopeType == "" || scopeID == "" {
		isSuperAdmin, err := as.repository.HasRole(actorUserID, domain.RoleSuperAdmin, domain.Scope{Type: domain.ScopePlatform, ID: domain.PlatformScopeID})
		if err != nil {
			return nil, err
		}
		if !isSuperAdmin {
			return nil, errors.New("not authorized to list grants")
		}
		return as.repository.ListGrants(userID, scopeType, scopeID)
	}

	if err := as.authorizeScopeRead(actorUserID, scopeType, scopeID); err != nil {
		return nil, err
	}

	return as.repository.ListGrants(userID, scopeType, scopeID)
}

func (as *AccessService) RecordResourceRelationship(resourceTypeValue, resourceID, relationValue, subjectTypeValue, subjectID string) error {
	relationship, err := parseRelationship(resourceTypeValue, resourceID, relationValue, subjectTypeValue, subjectID)
	if err != nil {
		return err
	}

	return as.repository.RecordResourceRelationship(relationship)
}

func (as *AccessService) authorizeGrant(actorUserID string, role domain.Role, scope domain.Scope) error {
	if actorUserID == "" {
		return errors.New("actor user id cannot be empty")
	}

	isSuperAdmin, err := as.repository.HasRole(actorUserID, domain.RoleSuperAdmin, domain.Scope{Type: domain.ScopePlatform, ID: domain.PlatformScopeID})
	if err != nil {
		return err
	}
	if isSuperAdmin {
		return nil
	}

	if role == domain.RoleSuperAdmin || role == domain.RoleAdmin {
		return errors.New("not authorized to grant role")
	}

	organizationID, err := as.organizationScopeID(scope)
	if err != nil {
		return err
	}

	isOwner, err := as.repository.HasOwnerRelationship(organizationID, actorUserID)
	if err != nil {
		return err
	}
	if isOwner {
		return nil
	}

	isAdmin, err := as.repository.HasRole(actorUserID, domain.RoleAdmin, domain.Scope{Type: domain.ScopeOrganization, ID: organizationID})
	if err != nil {
		return err
	}
	if isAdmin {
		return nil
	}

	return errors.New("not authorized to grant role")
}

func (as *AccessService) authorizeScopeRead(actorUserID, scopeTypeValue, scopeID string) error {
	scopeType, err := domain.ParseScopeType(scopeTypeValue)
	if err != nil {
		return err
	}

	scope, err := domain.NewScope(scopeType, scopeID)
	if err != nil {
		return err
	}

	isSuperAdmin, err := as.repository.HasRole(actorUserID, domain.RoleSuperAdmin, domain.Scope{Type: domain.ScopePlatform, ID: domain.PlatformScopeID})
	if err != nil {
		return err
	}
	if isSuperAdmin {
		return nil
	}

	organizationID, err := as.organizationScopeID(scope)
	if err != nil {
		return err
	}

	isOwner, err := as.repository.HasOwnerRelationship(organizationID, actorUserID)
	if err != nil {
		return err
	}
	if isOwner {
		return nil
	}

	isAdmin, err := as.repository.HasRole(actorUserID, domain.RoleAdmin, domain.Scope{Type: domain.ScopeOrganization, ID: organizationID})
	if err != nil {
		return err
	}
	if isAdmin {
		return nil
	}

	return errors.New("not authorized to view grants")
}

func (as *AccessService) organizationScopeID(scope domain.Scope) (string, error) {
	if scope.Type == domain.ScopeOrganization {
		return scope.ID, nil
	}

	resourceType := domain.ResourceType(scope.Type)
	resourceID := scope.ID

	for {
		parent, err := as.repository.FindParent(resourceType, resourceID)
		if err != nil {
			return "", err
		}
		if parent == nil {
			return "", errors.New("scope is not within an organization")
		}
		if parent.SubjectType == domain.SubjectOrganization {
			return parent.SubjectID, nil
		}

		resourceType = domain.ResourceType(parent.SubjectType)
		resourceID = parent.SubjectID
	}
}

func newAssignment(userID, roleValue, scopeTypeValue, scopeID string) (*domain.RoleAssignment, error) {
	role, err := domain.ParseRole(roleValue)
	if err != nil {
		return nil, err
	}

	scopeType, err := domain.ParseScopeType(scopeTypeValue)
	if err != nil {
		return nil, err
	}

	scope, err := domain.NewScope(scopeType, scopeID)
	if err != nil {
		return nil, err
	}

	return domain.NewRoleAssignment(userID, role, scope)
}

func parseRelationship(resourceTypeValue, resourceID, relationValue, subjectTypeValue, subjectID string) (*domain.ResourceRelationship, error) {
	return domain.NewResourceRelationship(
		domain.ResourceType(resourceTypeValue),
		resourceID,
		domain.Relation(relationValue),
		domain.SubjectType(subjectTypeValue),
		subjectID,
	)
}
