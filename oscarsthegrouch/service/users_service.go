package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"namaya/oscarsthegrouch/model"

	"github.com/google/uuid"
)

type UsersService interface {
	CreateUser(ctx context.Context, username string) (*model.User, error)
}

type usersService struct {
	dbClient *sql.DB
}

func NewUsersService(dbClient *sql.DB) UsersService {
	return &usersService{
		dbClient: dbClient,
	}
}

func (us *usersService) CreateUser(ctx context.Context, username string) (*model.User, error) {
	userId := uuid.New().String()

	result, err := us.dbClient.Exec("INSERT INTO users (id, name) VALUES (?, ?)", userId, username)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("CreateUser: no rows affected")
	}

	user := model.User{
		Id:   userId,
		Name: username,
	}

	log.Printf("User created with ID: %s", user.Id)

	return &user, nil
}
