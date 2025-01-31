package api

import (
	"encoding/json"
	"namaya/oscarsthegrouch/model"
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
	ctx := r.Context()

	var ballot model.Ballot

	if err := json.NewDecoder(r.Body).Decode(&ballot); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := b.ballotsService.SaveBallot(ctx, &ballot); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Ballot added successfully!"}`))
}

func (b *ballotEndpoint) getBallots(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "Hello, World 2!"}`))
}
