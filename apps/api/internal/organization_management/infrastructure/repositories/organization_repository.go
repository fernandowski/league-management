package repositories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"league-management/internal/organization_management/domain/organization"
	"league-management/internal/shared/app_errors"
	"league-management/internal/shared/database"
	"time"
)

type OrganizationRepository struct {
}

func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{}
}

func (or *OrganizationRepository) FindByName(orgName string, ownerID string) (*organization.Organization, error) {
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

	return &organization.Organization{Name: name, ID: &id, OrganizationOwnerId: ownerId}, nil
}

func (or *OrganizationRepository) FindById(orgId string) (*organization.Organization, error) {
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

	foundOrganization := organization.NewOrganization(&id, name, ownerId, true)
	return &foundOrganization, nil
}

func (ur *OrganizationRepository) Save(o *organization.Organization) error {
	pool := database.GetConnection()

	sql := "INSERT INTO league_management.organizations (name, user_id) VALUES ($1, $2)"

	_, err := pool.Exec(context.Background(), sql, o.Name, o.OrganizationOwnerId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return app_errors.ErrDuplicateResource
		}
	}

	return err
}

func (ur *OrganizationRepository) FetchAll(orgOwnerId string) ([]organization.Organization, error) {
	pool := database.GetConnection()
	sql := "SELECT id, name, user_id FROM league_management.organizations where user_id=$1"

	rows, err := pool.Query(context.Background(), sql, orgOwnerId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var organizations []organization.Organization

	for rows.Next() {
		var id, name, organizationOwnerId string
		var isActive bool

		if err := rows.Scan(&id, &name, &organizationOwnerId); err != nil {
			return nil, err
		}

		foundOrganization := organization.NewOrganization(&id, name, organizationOwnerId, isActive)
		organizations = append(organizations, foundOrganization)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return organizations, nil
}
