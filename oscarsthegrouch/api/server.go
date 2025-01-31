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

	// Build router
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	ballotsEndpoint := NewBallotsEndpoint(ballotsService)

	ballotsEndpoint.BuildRoutes(s)

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
