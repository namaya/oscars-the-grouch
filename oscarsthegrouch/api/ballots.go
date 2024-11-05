package api

import (
	"namaya/oscarsthegrouch/service"
	"net/http"

	"github.com/gorilla/mux"
)

type ballotEndpoint struct {
	ballotsService service.BallotsService
}

func NewBallotsEndpoint(bService service.BallotsService) Endpoint {
	return &ballotEndpoint{
		ballotsService: bService,
	}
}

func (b *ballotEndpoint) BuildRoutes(r *mux.Router) error {
	r.HandleFunc("/ballots", b.getBallots).Methods("GET")
	r.HandleFunc("/ballots", b.postBallot).Methods("POST")

	return nil
}

func (b *ballotEndpoint) postBallot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "Hello, World 1!"}`))
}

func (b *ballotEndpoint) getBallots(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "Hello, World 2!"}`))
}
