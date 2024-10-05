package main

import (
	routes "league-management/internal/user_management/interfaces/http"
	"log"
)

func main() {
	log.Print("initialized")
	routes.Initialize()
}
