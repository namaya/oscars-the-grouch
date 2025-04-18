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

// TODO: auth rights
func (e *gamesEndpoint) BuildRoutes(r *mux.Router) error {
	r.Handle("/games", e.RequireRightFunc(e.createGame)).Methods("POST")
	r.Handle("/games", e.RequireRightFunc(e.listGames)).Methods("GET")
	r.Handle("/games/{id}/players", e.RequireRightFunc(e.listPlayers)).Methods("GET")
	r.Handle("/games/{id}/scores", e.RequireRightFunc(e.scores)).Methods("GET")
	r.Handle("/games/{id}/ballots", e.RequireRightFunc(e.createBallot)).Methods("POST")
	r.Handle("/games/{id}/masterballot", e.RequireRightFunc(e.voteBallot)).Methods("PUT")

	return nil
}

type CreateGameRequest struct {
	Name string `json:"name"`
}

type CreateGameResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

func (b *gamesEndpoint) createGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createGameRequest CreateGameRequest

	if err := json.NewDecoder(r.Body).Decode(&createGameRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game, err := b.gamesService.CreateGame(ctx, createGameRequest.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resBody := CreateGameResponse{
		Id:    game.Id,
		Name:  game.Name,
		State: game.State,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resBody); err != nil {
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

type ListPlayersResponse struct {
	Players []PlayerResponse `json:"players"`
}

type PlayerResponse struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	AvatarUri string `json:"avatarUri"`
	Score     int    `json:"score"`
	State     string `json:"state"`
}

func (ge *gamesEndpoint) listPlayers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	players, err := ge.gamesService.ListPlayers(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	playersResp := make([]PlayerResponse, len(players))
	for i, player := range players {
		playersResp[i] = PlayerResponse{
			Id:        player.Id,
			Username:  player.User.Name,
			AvatarUri: player.User.AvatarUri,
			Score:     player.Score,
			State:     player.State,
		}
	}

	resBody := ListPlayersResponse{
		Players: playersResp,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ge *gamesEndpoint) scores(w http.ResponseWriter, r *http.Request) {
	// TODO: return the scores of all the players.
	// Requires that the game is active
}

func (ge *gamesEndpoint) createBallot(w http.ResponseWriter, r *http.Request) {
	// TODO: update player state to "Ready"

}

func (ge *gamesEndpoint) voteBallot(w http.ResponseWriter, r *http.Request) {
	// Requires that the game is active
}
