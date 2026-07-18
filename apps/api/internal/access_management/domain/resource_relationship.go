package domain

import (
	"errors"
	"strings"
)

type ResourceType string
type Relation string
type SubjectType string

const (
	ResourcePlatform     ResourceType = "platform"
	ResourceOrganization ResourceType = "organization"
	ResourceTeam         ResourceType = "team"
	ResourceLeague       ResourceType = "league"
	ResourceSeason       ResourceType = "season"
	ResourceMatch        ResourceType = "match"
)

const (
	RelationOwner  Relation = "owner"
	RelationParent Relation = "parent"
)

const (
	SubjectUser         SubjectType = "user"
	SubjectOrganization SubjectType = "organization"
	SubjectTeam         SubjectType = "team"
	SubjectLeague       SubjectType = "league"
	SubjectSeason       SubjectType = "season"
)

type ResourceRelationship struct {
	ResourceType ResourceType
	ResourceID   string
	Relation     Relation
	SubjectType  SubjectType
	SubjectID    string
}

func NewResourceRelationship(resourceType ResourceType, resourceID string, relation Relation, subjectType SubjectType, subjectID string) (*ResourceRelationship, error) {
	resourceID = strings.TrimSpace(resourceID)
	subjectID = strings.TrimSpace(subjectID)

	if !resourceType.IsValid() {
		return nil, errors.New("invalid resource type")
	}
	if resourceID == "" {
		return nil, errors.New("resource id cannot be empty")
	}
	if !relation.IsValid() {
		return nil, errors.New("invalid relation")
	}
	if !subjectType.IsValid() {
		return nil, errors.New("invalid subject type")
	}
	if subjectID == "" {
		return nil, errors.New("subject id cannot be empty")
	}

	return &ResourceRelationship{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Relation:     relation,
		SubjectType:  subjectType,
		SubjectID:    subjectID,
	}, nil
}

func (r ResourceType) IsValid() bool {
	switch r {
	case ResourcePlatform, ResourceOrganization, ResourceTeam, ResourceLeague, ResourceSeason, ResourceMatch:
		return true
	default:
		return false
	}
}

func (r Relation) IsValid() bool {
	switch r {
	case RelationOwner, RelationParent:
		return true
	default:
		return false
	}
}

func (s SubjectType) IsValid() bool {
	switch s {
	case SubjectUser, SubjectOrganization, SubjectTeam, SubjectLeague, SubjectSeason:
		return true
	default:
		return false
	}
}
