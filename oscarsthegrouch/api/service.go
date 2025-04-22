package api

import (
	"context"
	"net/http"
	"time"

	"namaya/oscarsthegrouch/database"
	"namaya/oscarsthegrouch/log"
	"namaya/oscarsthegrouch/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Endpoint interface {
	BuildRoutes(r *mux.Router) error
}

func ServerHandler() {
	ctx := log.WithFields(context.Background())
	logger := log.Get(ctx)

	// Build infrastructure
	dbClient, err := database.ConnectDb(ctx)
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}

	// Build services
	usersService := service.NewUsersService(dbClient)
	gamesService := service.NewGameService(dbClient)
	ballotsService := service.NewBallotsService(dbClient)

	// Build API endpoints
	ae := NewAuthorizedEndpoint(usersService)

	gamesEndpoint := NewGamesEndpoint(ae, gamesService, usersService)
	ballotsEndpoint := NewBallotsEndpoint(ballotsService)
	usersEndpoint := NewUsersEndpoint(usersService)

	r, err := BuildRouter(gamesEndpoint, usersEndpoint, ballotsEndpoint)
	if err != nil {
		logger.Fatalf("Error building router: %v", err)
	}

	// Build static file server
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Start server
	logger.Info("Starting server on :8080")

	http.ListenAndServe(":8080", r)
}

func BuildRouter(endpoints ...Endpoint) (*mux.Router, error) {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	r.Use(Trace())
	r.Use(ResponseLogger())

	for _, e := range endpoints {
		if err := e.BuildRoutes(s); err != nil {
			return nil, err
		}
	}

	return r, nil
}

type statusCapture struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusCapture) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func (sc *statusCapture) Write(b []byte) (int, error) {
	return sc.ResponseWriter.Write(b)
}

func ResponseLogger() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sc := &statusCapture{w, http.StatusOK}

			next.ServeHTTP(sc, r)

			dt := time.Since(start)
			logger := log.Get(r.Context()).With(
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sc.statusCode),
				zap.Duration("duration", dt),
			)

			if sc.statusCode >= 500 {
				logger.Errorf("5xx error for request (%s)", dt)
			} else if sc.statusCode >= 400 {
				logger.Warnf("4xx error for request (%s)", dt)
			} else {
				logger.Infof("request (%s)", dt)
			}
		})
	}
}

func Trace() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			traceId := r.Header.Get("X-Trace-Id")
			if traceId == "" {
				traceId = uuid.New().String()
			}

			ctx = log.WithFields(ctx, zap.String("traceId", traceId))
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
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

		ctx = log.WithFields(ctx, zap.String("userId", userId))
		ctx = context.WithValue(ctx, "userId", userId)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
