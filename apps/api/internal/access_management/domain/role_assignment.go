package domain

import "errors"

type RoleAssignment struct {
	ID     string
	UserID string
	Role   Role
	Scope  Scope
}

func NewRoleAssignment(userID string, role Role, scope Scope) (*RoleAssignment, error) {
	if userID == "" {
		return nil, errors.New("user id cannot be empty")
	}

	if !role.IsValid() {
		return nil, errors.New("invalid role")
	}

	if err := validateRoleScope(role, scope.Type); err != nil {
		return nil, err
	}

	return &RoleAssignment{UserID: userID, Role: role, Scope: scope}, nil
}

func validateRoleScope(role Role, scopeType ScopeType) error {
	switch role {
	case RoleSuperAdmin:
		if scopeType != ScopePlatform {
			return errors.New("super_admin must be platform scoped")
		}
	case RoleAdmin:
		if scopeType != ScopeOrganization {
			return errors.New("admin must be organization scoped")
		}
	case RoleCoach, RoleTeamMember:
		if scopeType != ScopeTeam {
			return errors.New("team role must be team scoped")
		}
	case RoleReferee:
		if scopeType != ScopeLeague && scopeType != ScopeSeason && scopeType != ScopeMatch {
			return errors.New("referee must be league, season, or match scoped")
		}
	default:
		return errors.New("invalid role")
	}

	return nil
}
