package domain

import (
	"errors"
	"strings"
)

type ScopeType string

const (
	ScopePlatform     ScopeType = "platform"
	ScopeOrganization ScopeType = "organization"
	ScopeTeam         ScopeType = "team"
	ScopeLeague       ScopeType = "league"
	ScopeSeason       ScopeType = "season"
	ScopeMatch        ScopeType = "match"
)

const PlatformScopeID = "platform"

func ParseScopeType(value string) (ScopeType, error) {
	scopeType := ScopeType(strings.TrimSpace(value))
	if scopeType.IsValid() {
		return scopeType, nil
	}

	return "", errors.New("invalid scope type")
}

func (s ScopeType) IsValid() bool {
	switch s {
	case ScopePlatform, ScopeOrganization, ScopeTeam, ScopeLeague, ScopeSeason, ScopeMatch:
		return true
	default:
		return false
	}
}

type Scope struct {
	Type ScopeType
	ID   string
}

func NewScope(scopeType ScopeType, scopeID string) (Scope, error) {
	if !scopeType.IsValid() {
		return Scope{}, errors.New("invalid scope type")
	}

	scopeID = strings.TrimSpace(scopeID)
	if scopeType == ScopePlatform {
		scopeID = PlatformScopeID
	}

	if scopeID == "" {
		return Scope{}, errors.New("scope id cannot be empty")
	}

	return Scope{Type: scopeType, ID: scopeID}, nil
}
