package services

import (
	"league-management/internal/access_management/domain"
	"testing"
)

type fakeAccessRepo struct {
	grants        map[string]bool
	relationships map[string]*domain.ResourceRelationship
}

func newFakeAccessRepo() *fakeAccessRepo {
	return &fakeAccessRepo{
		grants:        map[string]bool{},
		relationships: map[string]*domain.ResourceRelationship{},
	}
}

func (r *fakeAccessRepo) GrantRole(assignment *domain.RoleAssignment) error {
	r.grants[grantKey(assignment.UserID, assignment.Role, assignment.Scope)] = true
	return nil
}

func (r *fakeAccessRepo) RevokeRole(assignment *domain.RoleAssignment) error {
	delete(r.grants, grantKey(assignment.UserID, assignment.Role, assignment.Scope))
	return nil
}

func (r *fakeAccessRepo) HasRole(userID string, role domain.Role, scope domain.Scope) (bool, error) {
	return r.grants[grantKey(userID, role, scope)], nil
}

func (r *fakeAccessRepo) ListGrants(_, _, _ string) ([]domain.RoleAssignment, error) {
	return nil, nil
}

func (r *fakeAccessRepo) RecordResourceRelationship(relationship *domain.ResourceRelationship) error {
	r.relationships[relationshipKey(relationship.ResourceType, relationship.ResourceID, relationship.Relation)] = relationship
	return nil
}

func (r *fakeAccessRepo) FindParent(resourceType domain.ResourceType, resourceID string) (*domain.ResourceRelationship, error) {
	relationship, ok := r.relationships[relationshipKey(resourceType, resourceID, domain.RelationParent)]
	if !ok {
		return nil, nil
	}
	return relationship, nil
}

func (r *fakeAccessRepo) HasOwnerRelationship(organizationID, userID string) (bool, error) {
	relationship, ok := r.relationships[relationshipKey(domain.ResourceOrganization, organizationID, domain.RelationOwner)]
	return ok && relationship.SubjectType == domain.SubjectUser && relationship.SubjectID == userID, nil
}

func TestAccessService_OrgOwnerCanGrantTeamRoleForTeamInOrganization(t *testing.T) {
	repo := newFakeAccessRepo()
	service := NewAccessService(repo)

	mustRecord(t, service, "organization", "org-1", "owner", "user", "owner-1")
	mustRecord(t, service, "team", "team-1", "parent", "organization", "org-1")

	err := service.GrantRole("owner-1", "coach-1", "coach", "team", "team-1")
	if err != nil {
		t.Fatalf("expected grant to succeed, got %v", err)
	}

	hasRole, err := service.HasRole("coach-1", "coach", "team", "team-1")
	if err != nil {
		t.Fatalf("expected has role to succeed, got %v", err)
	}
	if !hasRole {
		t.Fatalf("expected user to have coach role")
	}
}

func TestAccessService_OrgAdminCanGrantTeamMemberForTeamInOrganization(t *testing.T) {
	repo := newFakeAccessRepo()
	service := NewAccessService(repo)

	mustGrantDirect(t, repo, "admin-1", domain.RoleAdmin, domain.Scope{Type: domain.ScopeOrganization, ID: "org-1"})
	mustRecord(t, service, "team", "team-1", "parent", "organization", "org-1")

	err := service.GrantRole("admin-1", "member-1", "team_member", "team", "team-1")
	if err != nil {
		t.Fatalf("expected grant to succeed, got %v", err)
	}
}

func TestAccessService_UserCannotGrantTeamRoleOutsideTheirOrganization(t *testing.T) {
	repo := newFakeAccessRepo()
	service := NewAccessService(repo)

	mustRecord(t, service, "organization", "org-1", "owner", "user", "owner-1")
	mustRecord(t, service, "team", "team-2", "parent", "organization", "org-2")

	err := service.GrantRole("owner-1", "coach-1", "coach", "team", "team-2")
	if err == nil {
		t.Fatalf("expected grant to fail")
	}
}

func TestAccessService_OrgOwnerCannotGrantSuperAdmin(t *testing.T) {
	repo := newFakeAccessRepo()
	service := NewAccessService(repo)

	mustRecord(t, service, "organization", "org-1", "owner", "user", "owner-1")

	err := service.GrantRole("owner-1", "admin-1", "super_admin", "platform", "platform")
	if err == nil {
		t.Fatalf("expected grant to fail")
	}
}

func TestAccessService_SuperAdminCanGrantAnyRole(t *testing.T) {
	repo := newFakeAccessRepo()
	service := NewAccessService(repo)

	mustGrantDirect(t, repo, "super-1", domain.RoleSuperAdmin, domain.Scope{Type: domain.ScopePlatform, ID: domain.PlatformScopeID})

	err := service.GrantRole("super-1", "admin-1", "admin", "organization", "org-1")
	if err != nil {
		t.Fatalf("expected grant to succeed, got %v", err)
	}
}

func TestAccessService_DelegatesOrganizationAdministration(t *testing.T) {
	repo := newFakeAccessRepo()
	service := NewAccessService(repo)

	mustRecord(t, service, "team", "team-in-org-1", "parent", "organization", "org-1")
	mustRecord(t, service, "team", "team-in-org-2", "parent", "organization", "org-2")

	if err := service.BootstrapSuperAdmin("super-admin"); err != nil {
		t.Fatalf("expected super-admin bootstrap to succeed, got %v", err)
	}

	t.Run("super admin grants organization admin", func(t *testing.T) {
		err := service.GrantRole("super-admin", "org-admin", "admin", "organization", "org-1")
		if err != nil {
			t.Fatalf("expected admin grant to succeed, got %v", err)
		}
	})

	t.Run("organization admin grants a team coach in their organization", func(t *testing.T) {
		err := service.GrantRole("org-admin", "coach", "coach", "team", "team-in-org-1")
		if err != nil {
			t.Fatalf("expected coach grant to succeed, got %v", err)
		}

		hasRole, err := service.HasRole("coach", "coach", "team", "team-in-org-1")
		if err != nil {
			t.Fatalf("expected role check to succeed, got %v", err)
		}
		if !hasRole {
			t.Fatalf("expected coach role to be granted")
		}
	})

	t.Run("organization admin cannot grant access outside their organization", func(t *testing.T) {
		err := service.GrantRole("org-admin", "other-coach", "coach", "team", "team-in-org-2")
		if err == nil {
			t.Fatalf("expected cross-organization grant to fail")
		}
	})
}

func TestAccessService_BootstrapSuperAdminGrantsPlatformRole(t *testing.T) {
	repo := newFakeAccessRepo()
	service := NewAccessService(repo)

	if err := service.BootstrapSuperAdmin("super-1"); err != nil {
		t.Fatalf("expected bootstrap to succeed, got %v", err)
	}

	hasRole, err := service.HasRole("super-1", "super_admin", "platform", "platform")
	if err != nil {
		t.Fatalf("expected role check to succeed, got %v", err)
	}
	if !hasRole {
		t.Fatalf("expected user to have super_admin role")
	}
}

func TestAccessService_RejectsInvalidRoleScope(t *testing.T) {
	repo := newFakeAccessRepo()
	service := NewAccessService(repo)

	err := service.GrantRole("super-1", "coach-1", "coach", "organization", "org-1")
	if err == nil {
		t.Fatalf("expected invalid role scope error")
	}
}

func mustRecord(t *testing.T, service *AccessService, resourceType, resourceID, relation, subjectType, subjectID string) {
	t.Helper()
	if err := service.RecordResourceRelationship(resourceType, resourceID, relation, subjectType, subjectID); err != nil {
		t.Fatalf("expected relationship record to succeed, got %v", err)
	}
}

func mustGrantDirect(t *testing.T, repo *fakeAccessRepo, userID string, role domain.Role, scope domain.Scope) {
	t.Helper()
	assignment, err := domain.NewRoleAssignment(userID, role, scope)
	if err != nil {
		t.Fatalf("expected assignment to be valid, got %v", err)
	}
	if err := repo.GrantRole(assignment); err != nil {
		t.Fatalf("expected grant to succeed, got %v", err)
	}
}

func grantKey(userID string, role domain.Role, scope domain.Scope) string {
	return userID + "|" + string(role) + "|" + string(scope.Type) + "|" + scope.ID
}

func relationshipKey(resourceType domain.ResourceType, resourceID string, relation domain.Relation) string {
	return string(resourceType) + "|" + resourceID + "|" + string(relation)
}

var _ accessRepository = (*fakeAccessRepo)(nil)
