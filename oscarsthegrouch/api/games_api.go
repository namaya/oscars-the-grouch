package api

import (
	"encoding/json"
	// "namaya/oscarsthegrouch/model"
	"namaya/oscarsthegrouch/model"
	"namaya/oscarsthegrouch/service"
	"net/http"

	"github.com/gorilla/mux"
)

type gamesEndpoint struct {
	AuthorizedEndpoint
	gamesService service.GamesService
	usersService service.UsersService
}

func NewGamesEndpoint(ae AuthorizedEndpoint, gs service.GamesService, us service.UsersService) Endpoint {
	return &gamesEndpoint{
		AuthorizedEndpoint: ae,
		gamesService:       gs,
		usersService:       us,
	}
}

// TODO: auth rights
func (e *gamesEndpoint) BuildRoutes(r *mux.Router) error {
	r.Handle("/games", e.RequireRightFunc(e.createGame)).Methods("POST")
	r.Handle("/games", e.RequireRightFunc(e.listGames)).Methods("GET")
	r.Handle("/games/{id}/players", e.RequireRightFunc(e.listPlayers)).Methods("GET")
	r.Handle("/games/{id}/players", e.RequireRightFunc(e.addPlayer)).Methods("POST")
	r.Handle("/games/{gid}/players/{pid}/ballots", e.RequireRightFunc(e.createBallot)).Methods("POST")
	r.Handle("/games/{id}/scores", e.RequireRightFunc(e.scores)).Methods("GET")
	r.Handle("/games/{id}/masterballot", e.RequireRightFunc(e.voteBallot)).Methods("PUT")

	r.Handle("/games/{id}/nominations", e.RequireRightFunc(e.getNominations)).Methods("GET")

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
	UserId    string `json:"userId"`
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
			UserId:    player.User.Id,
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

type AddPlayerRequest struct {
	Name      string `json:"name"`
	AvatarUri string `json:"avatarUri"`
}

func (ge *gamesEndpoint) addPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	gameId := vars["id"]
	if gameId == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	var addPlayerRequest AddPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&addPlayerRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if addPlayerRequest.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if addPlayerRequest.AvatarUri == "" {
		http.Error(w, "Avatar URI is required", http.StatusBadRequest)
		return
	}

	// Create a new user
	user, err := ge.usersService.CreateUser(ctx, addPlayerRequest.Name, addPlayerRequest.AvatarUri)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err := ge.gamesService.AddPlayer(ctx, gameId, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pr := PlayerResponse{
		Id:        player.Id,
		UserId:    player.User.Id,
		Username:  player.User.Name,
		AvatarUri: player.User.AvatarUri,
		Score:     player.Score,
		State:     player.State,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (ge *gamesEndpoint) getNominations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nominations, err := ge.gamesService.GetNominations(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(nominations); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ge *gamesEndpoint) scores(w http.ResponseWriter, r *http.Request) {
	// TODO: return the scores of all the players.
	// Requires that the game is active
}

type CreateBallotRequest struct {
	Votes []*model.Vote `json:"votes"`
}

func (ge *gamesEndpoint) createBallot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	gid := vars["gid"]
	if gid == "" {
		http.Error(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	pid := vars["pid"]
	if pid == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}

	var createBallotRequest model.Ballot
	if err := json.NewDecoder(r.Body).Decode(&createBallotRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ballot, err := ge.gamesService.CreateBallot(ctx, gid, pid, createBallotRequest.Votes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ge.gamesService.UpdatePlayerState(ctx, gid, ballot.PlayerId, "Ready")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ballot); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ge *gamesEndpoint) voteBallot(w http.ResponseWriter, r *http.Request) {
	// Requires that the game is active
}
