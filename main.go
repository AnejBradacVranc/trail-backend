package main

import (
	"backend/internal/handlers"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	r := http.NewServeMux()
	handlers.Handler(r)
	fmt.Println("Starting server on port 8080")

	withCors := handlers.EnableCORS(r)

	err := http.ListenAndServe(":8080", withCors)
	if err != nil {
		log.Error(err)
		return
	}
}
