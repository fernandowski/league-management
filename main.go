package main

import (
	"league-management/internal/user_management/interfaces/http"
	"log"
)

func main() {
	log.Print("initialized")
	routes.Initialize()
}
