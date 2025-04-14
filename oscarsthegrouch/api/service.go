package api

import (
	"context"
	"log"
	"net/http"

	"namaya/oscarsthegrouch/database"
	"namaya/oscarsthegrouch/service"

	"github.com/gorilla/mux"
)

type Endpoint interface {
	BuildRoutes(r *mux.Router) error
}

func ServerHandler() {
	// Build infrastructure
	ctx := context.Background()

	dbClient, err := database.ConnectDb(ctx)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Build services
	ballotsService := service.NewBallotsService(dbClient)
	usersService := service.NewUsersService(dbClient)

	// Build endpoints
	gamesEndpoint := NewGamesEndpoint()
	ballotsEndpoint := NewBallotsEndpoint(ballotsService)
	usersEndpoint := NewUsersEndpoint(usersService)

	r, err := BuildRouter(gamesEndpoint, usersEndpoint, ballotsEndpoint)
	if err != nil {
		log.Fatalf("Error building router: %v", err)
	}

	http.Handle("/", r)

	// Start server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

func BuildRouter(endpoints ...Endpoint) (*mux.Router, error) {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	for _, e := range endpoints {
		if err := e.BuildRoutes(s); err != nil {
			return nil, err
		}
	}

	return r, nil
}
