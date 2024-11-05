package api

import (
	"namaya/oscarsthegrouch/database"
	"namaya/oscarsthegrouch/service"
	"net/http"

	"github.com/gorilla/mux"
)

type Endpoint interface {
	BuildRoutes(r *mux.Router) error
}

func ServerHandler() {
	// Build infrastructure
	err := database.RunMigrations()
	if err != nil {
		panic(err)
	}

	// Build services
	ballotsService := service.NewBallotsService()

	// Build router
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	ballotsEndpoint := NewBallotsEndpoint(ballotsService)

	ballotsEndpoint.BuildRoutes(s)

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
