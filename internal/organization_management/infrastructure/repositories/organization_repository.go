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

	log.Print(orgName, ownerID)

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

	return &domain.Organization{Name: name, ID: id, OrganizationOwnerId: ownerId}, nil
}

func (or *OrganizationRepository) FindById(orgId string) (*domain.Organization, error) {
	connection := database.GetConnection()

	var id string
	var name string
	var dateCreated time.Time
	var dateUpdated time.Time
	var ownerId string

	sql := "SELECT id, name, user_id, created_at, updated_at FROM league_management.organizations where id=$1"

	err := connection.QueryRow(context.Background(), sql, orgId).Scan(id, name, ownerId, dateCreated, dateUpdated)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		panic(err)
	}

	return &domain.Organization{Name: name, ID: id, OrganizationOwnerId: ownerId}, nil
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
