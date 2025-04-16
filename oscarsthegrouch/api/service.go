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
	ae := NewAuthorizedEndpoint(usersService)

	gamesEndpoint := NewGamesEndpoint(ae, gamesService)
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

type AuthorizedEndpoint interface {
	RequireRightFunc(next http.HandlerFunc, rights ...string) http.Handler
}

type authorizedEndpoint struct {
	usersService service.UsersService
}

func NewAuthorizedEndpoint(us service.UsersService) AuthorizedEndpoint {
	return &authorizedEndpoint{
		usersService: us,
	}
}

func (ae *authorizedEndpoint) RequireRightFunc(next http.HandlerFunc, rights ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("Authorization")
		if userId == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		_, err := ae.usersService.GetUser(ctx, userId)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(ctx, "userId", userId))

		next.ServeHTTP(w, r)
	})
}
