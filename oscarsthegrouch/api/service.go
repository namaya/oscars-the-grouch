package api

import (
	"context"
	"log"
	"net/http"
	"os"

	"namaya/oscarsthegrouch/database"
	"namaya/oscarsthegrouch/service"

	"github.com/gorilla/handlers"
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
	usersService := service.NewUsersService(dbClient)
	gamesService := service.NewGameService(dbClient)
	ballotsService := service.NewBallotsService(dbClient)

	// Build API endpoints
	gamesEndpoint := NewGamesEndpoint(gamesService)
	ballotsEndpoint := NewBallotsEndpoint(ballotsService)
	usersEndpoint := NewUsersEndpoint(usersService)

	r, err := BuildRouter(gamesEndpoint, usersEndpoint, ballotsEndpoint)
	if err != nil {
		log.Fatalf("Error building router: %v", err)
	}

	// Build static file server
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Start server
	log.Println("Starting server on :8080")

	logr := handlers.LoggingHandler(os.Stdout, r)

	http.ListenAndServe(":8080", logr)
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

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("%s")
	})
}
