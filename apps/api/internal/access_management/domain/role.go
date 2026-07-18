package domain

import (
	"errors"
	"strings"
)

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleCoach      Role = "coach"
	RoleTeamMember Role = "team_member"
	RoleReferee    Role = "referee"
)

func ParseRole(value string) (Role, error) {
	role := Role(strings.TrimSpace(value))
	if role.IsValid() {
		return role, nil
	}

	return "", errors.New("invalid role")
}

func (r Role) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RoleCoach, RoleTeamMember, RoleReferee:
		return true
	default:
		return false
	}
}
