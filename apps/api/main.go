package main

import (
	apphttp "league-management/internal/app/http"
	"log"
)

func main() {
	log.Print("initialized")
	if err := apphttp.Run(); err != nil {
		log.Fatal(err)
	}
}
