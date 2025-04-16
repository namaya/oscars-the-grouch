package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"namaya/oscarsthegrouch/model"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type UsersService interface {
	CreateUser(ctx context.Context, username string, avatar string) (*model.User, error)
	GetUser(ctx context.Context, userId string) (*model.User, error)
	ListAvatars(ctx context.Context) ([]string, error)
}

type usersService struct {
	dbClient *sql.DB
}

func NewUsersService(dbClient *sql.DB) UsersService {
	return &usersService{
		dbClient: dbClient,
	}
}

func (us *usersService) CreateUser(ctx context.Context, username string, avatarUri string) (*model.User, error) {
	userId := uuid.New().String()

	result, err := us.dbClient.Exec("INSERT INTO users (id, name, avatar_uri) VALUES (?, ?, ?)", userId, username, avatarUri)
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
		Id:        userId,
		Name:      username,
		AvatarUri: avatarUri,
	}

	log.Printf("User created with ID: %s", user.Id)

	return &user, nil
}

func (us *usersService) ListAvatars(ctx context.Context) ([]string, error) {
	files := []string{}

	err := filepath.Walk("./static/avatars", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Walk: %w", err)
		}

		if !info.IsDir() {
			path := "/static/avatars/" + info.Name()
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ListAvatars: %w", err)
	}

	return files, nil
}

func (us *usersService) GetUser(ctx context.Context, userId string) (*model.User, error) {
	row := us.dbClient.QueryRow("SELECT id, name, avatar_uri FROM users WHERE id = ?", userId)

	user := &model.User{}

	err := row.Scan(&user.Id, &user.Name, &user.AvatarUri)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("GetUser: %w", err)
	}

	return user, nil
}
