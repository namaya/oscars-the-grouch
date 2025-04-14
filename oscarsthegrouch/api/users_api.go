package api

import (
	"encoding/json"
	"namaya/oscarsthegrouch/model"
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
	r.HandleFunc("/users", ue.postUser).Methods("POST")

	return nil
}

// postUser handles HTTP requests to create a new user.
func (ue *usersEndpoint) postUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ue.usersService.SaveUser(ctx, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := PostUserResponse{
		Message: "New user created.",
		User:    user,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type PostUserResponse struct {
	Message string     `json:"message"`
	User    model.User `json:"user"`
}
