package api

import (
	"encoding/json"
	"namaya/oscarsthegrouch/service"
	"net/http"

	"github.com/gorilla/mux"
)

type usersEndpoint struct {
	usersService service.UsersService
}

func NewUsersEndpoint(us service.UsersService) Endpoint {
	return &usersEndpoint{
		usersService: us,
	}
}

func (ue *usersEndpoint) BuildRoutes(r *mux.Router) error {
	r.HandleFunc("/users", ue.createUser).Methods("POST")
	r.HandleFunc("/users/avatars", ue.listAvatars).Methods("GET")
	r.HandleFunc("/users/{id}", ue.getUser).Methods("GET")

	return nil
}

type CreateUserRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type CreateUserResponse struct {
	UserId string `json:"userId"`
}

// createUser handles HTTP requests to create a new user.
func (ue *usersEndpoint) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ue.usersService.CreateUser(ctx, reqBody.Name, reqBody.Avatar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resBody := CreateUserResponse{
		UserId: user.Id,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ListAvatarsResponse struct {
	Avatars []string `json:"avatars"`
}

func (ue *usersEndpoint) listAvatars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	avatars, err := ue.usersService.ListAvatars(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resBody := ListAvatarsResponse{
		Avatars: avatars,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type GetUserResponse struct {
	UserId    string `json:"userId"`
	Name      string `json:"name"`
	AvatarUri string `json:"avatarUri"`
}

func (ue *usersEndpoint) getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := ue.usersService.GetUser(ctx, id)
	if err != nil {
		if err == service.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resBody := GetUserResponse{
		UserId:    user.Id,
		Name:      user.Name,
		AvatarUri: user.AvatarUri,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
