package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	accessservices "league-management/internal/access_management/application/services"
	accesspg "league-management/internal/access_management/infrastructure/repositories/postgres"
	userpg "league-management/internal/user_management/infrastructure/repositories/postgres"
)

func main() {
	email := flag.String("email", "", "email address of the existing user to grant super_admin")
	flag.Parse()

	if strings.TrimSpace(*email) == "" {
		log.Fatal("email is required: --email=user@example.com")
	}

	user := userpg.NewUserRepository().FindByEmail(*email)
	if user == nil {
		log.Fatalf("user with email %q does not exist", *email)
	}

	service := accessservices.NewAccessService(accesspg.NewAccessRepository())
	if err := service.BootstrapSuperAdmin(user.Id); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("granted super_admin to %s (%s)\n", user.Email, user.Id)
}
