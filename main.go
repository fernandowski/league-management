package main

import (
	routes "league-management/internal/interfaces/http"
	"log"
)

func main() {
	log.Print("initialized")
	routes.Initialize()
}
