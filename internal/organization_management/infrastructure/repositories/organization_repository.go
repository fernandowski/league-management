package repositories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"league-management/internal/organization_management/domain"
	"league-management/internal/shared/database"
	"log"
	"time"
)

type OrganizationRepository struct {
}

func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{}
}

func (or *OrganizationRepository) FindByName(orgName string, ownerID string) (*domain.Organization, error) {
	connection := database.GetConnection()

	sql := "SELECT id, name, user_id, created_at, updated_at FROM league_management.organizations where name=$1 AND user_id=$2"

	var id string
	var name string
	var dateCreated time.Time
	var dateUpdated time.Time
	var ownerId string

	err := connection.QueryRow(context.Background(), sql, orgName, ownerID).Scan(&id, &name, &ownerId, &dateCreated, &dateUpdated)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		panic(err)
	}

	return &domain.Organization{Name: name, ID: &id, OrganizationOwnerId: ownerId}, nil
}

func (or *OrganizationRepository) FindById(orgId string) (*domain.Organization, error) {
	connection := database.GetConnection()

	var id string
	var name string
	var dateCreated time.Time
	var dateUpdated time.Time
	var ownerId string

	sql := "SELECT id, name, user_id, created_at, updated_at FROM league_management.organizations where id=$1"

	err := connection.QueryRow(context.Background(), sql, orgId).Scan(&id, &name, &ownerId, &dateCreated, &dateUpdated)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		panic(err)
	}

	organization := domain.NewOrganization(&id, name, ownerId, true)
	return &organization, nil
}

func (ur *OrganizationRepository) Save(o *domain.Organization) error {
	pool := database.GetConnection()

	sql := "INSERT INTO league_management.organizations (name, user_id) VALUES ($1, $2)"

	_, err := pool.Exec(context.Background(), sql, o.Name, o.OrganizationOwnerId)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (ur *OrganizationRepository) FetchAll(orgOwnerId string) ([]domain.Organization, error) {
	pool := database.GetConnection()
	sql := "SELECT id, name, user_id FROM league_management.organizations where user_id=$1"

	rows, err := pool.Query(context.Background(), sql, orgOwnerId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var organizations []domain.Organization

	for rows.Next() {
		var id, name, organizationOwnerId string
		var isActive bool

		if err := rows.Scan(&id, &name, &organizationOwnerId); err != nil {
			return nil, err
		}

		organization := domain.NewOrganization(&id, name, organizationOwnerId, isActive)
		organizations = append(organizations, organization)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return organizations, nil
}
