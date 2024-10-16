package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ServerHandler() {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	s.HandleFunc("/ballots", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{"message": "Hello, World!"}`))
	}).Methods("POST")

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
