package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Endpoint interface {
	BuildRoutes(r *mux.Router) error
}

func ServerHandler() {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	ballotsEndpoint := NewBallotsEndpoint()

	ballotsEndpoint.BuildRoutes(s)

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
