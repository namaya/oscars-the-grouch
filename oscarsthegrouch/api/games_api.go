package api

import (
	"encoding/json"
	// "namaya/oscarsthegrouch/model"
	"namaya/oscarsthegrouch/service"
	"net/http"

	"github.com/gorilla/mux"
)

type gamesEndpoint struct {
	AuthorizedEndpoint
	gamesService service.GamesService
}

func NewGamesEndpoint(ae AuthorizedEndpoint, gs service.GamesService) Endpoint {
	return &gamesEndpoint{
		AuthorizedEndpoint: ae,
		gamesService:       gs,
	}
}

func (e *gamesEndpoint) BuildRoutes(r *mux.Router) error {
	r.Handle("/games", e.RequireRightFunc(e.createGame)).Methods("POST")
	r.Handle("/games", e.RequireRightFunc(e.listGames)).Methods("GET")

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
	ctx := r.Context()

	games, err := ge.gamesService.ListGames(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gamesResp := make([]GameResponse, len(games))
	for i, game := range games {
		gamesResp[i] = GameResponse{
			Id:    game.Id,
			Name:  game.Name,
			State: game.State,
		}
	}

	resBody := ListGamesResponse{
		Games: gamesResp,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
