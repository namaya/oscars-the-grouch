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

	return nil
}

// createUser handles HTTP requests to create a new user.
func (ue *usersEndpoint) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ue.usersService.CreateUser(ctx, reqBody.Name)
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

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	UserId string `json:"userId"`
}
