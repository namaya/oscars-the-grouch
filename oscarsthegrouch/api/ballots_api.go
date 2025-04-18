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

func NewBallotsEndpoint(bs service.BallotsService) Endpoint {
	return &ballotEndpoint{
		ballotsService: bs,
	}
}

func (e *ballotEndpoint) BuildRoutes(r *mux.Router) error {
	r.HandleFunc("/ballots", e.getBallots).Methods("GET")
	r.HandleFunc("/ballots", e.createBallot).Methods("POST")
	r.HandleFunc("/ballots/{id}", e.updateBallot).Methods("PATCH")

	return nil
}

func (e *ballotEndpoint) createBallot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var ballot model.Ballot

	if err := json.NewDecoder(r.Body).Decode(&ballot); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := e.ballotsService.SaveBallot(ctx, &ballot); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Ballot added successfully!"}`))
}

func (b *ballotEndpoint) getBallots(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "Hello, World 2!"}`))
}

type PatchBallotRequest struct {
	CategoryId string `json:"categoryId"`
	Vote       string `json:"vote"`
}

func (b *ballotEndpoint) updateBallot(w http.ResponseWriter, r *http.Request) {
	// TODO: update a single vote

}
