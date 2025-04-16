package api

import (
	"encoding/json"
	// "namaya/oscarsthegrouch/model"
	"namaya/oscarsthegrouch/service"
	"net/http"

	"github.com/gorilla/mux"
)

type gamesEndpoint struct {
	gamesService service.GamesService
}

func NewGamesEndpoint(gs service.GamesService) Endpoint {
	return &gamesEndpoint{
		gamesService: gs,
	}
}

func (b *gamesEndpoint) BuildRoutes(r *mux.Router) error {
	r.HandleFunc("/games", b.createGame).Methods("POST")
	r.HandleFunc("/games", b.listGames).Methods("GET")

	return nil
}

type CreateGameRequest struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

type CreateGameResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

func (b *gamesEndpoint) createGame(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	userId := r.Header.Get("Authorization")
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var createGameRequest CreateGameRequest

	if err := json.NewDecoder(r.Body).Decode(&createGameRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := CreateGameResponse{
		Name:  createGameRequest.Name,
		State: createGameRequest.State,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ListGamesResponse struct {
	Games []GameResponse `json:"games"`
}

type GameResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

func (ge *gamesEndpoint) listGames(w http.ResponseWriter, r *http.Request) {
	// TODO: check that user exists
	userId := r.Header.Get("Authorization")
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	games, err := ge.gamesService.ListGames(ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resBody := make([]GameResponse, len(games))
	for i, game := range games {
		resBody[i] = GameResponse{
			Id:    game.Id,
			Name:  game.Name,
			State: game.State,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
