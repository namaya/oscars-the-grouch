package service

import (
	"context"
	"database/sql"
	"fmt"
	"namaya/oscarsthegrouch/model"

	"github.com/google/uuid"
)

type GamesService interface {
	ListGames(ctx context.Context) ([]*model.Game, error)
	CreateGame(ctx context.Context, name string) (*model.Game, error)
}

type gamesService struct {
	dbClient *sql.DB
}

func NewGameService(dbClient *sql.DB) GamesService {
	return &gamesService{
		dbClient: dbClient,
	}
}

func (gs *gamesService) ListGames(ctx context.Context) ([]*model.Game, error) {
	userId := ctx.Value("userId").(string)

	rows, err := gs.dbClient.QueryContext(ctx, "SELECT id, name, state FROM games WHERE owner_id = ?", userId)
	if err != nil {
		return nil, fmt.Errorf("ListGames: %w", err)
	}
	defer rows.Close()

	var games []*model.Game
	for rows.Next() {
		var game model.Game
		if err := rows.Scan(&game.Id, &game.Name, &game.State); err != nil {
			return nil, fmt.Errorf("ListGames: %w", err)
		}
		games = append(games, &game)
	}

	return games, nil
}

func (gs *gamesService) CreateGame(ctx context.Context, name string) (*model.Game, error) {
	userId := ctx.Value("userId").(string)
	gameId := uuid.New().String()
	state := "Created"

	result, err := gs.dbClient.ExecContext(ctx, "INSERT INTO games (id, name, state, owner_id) VALUES (?, ?, ?, ?)", gameId, name, state, userId)
	if err != nil {
		return nil, fmt.Errorf("CreateGame: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("CreateGame: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("CreateGame: no rows affected")
	}

	game := &model.Game{
		Id:    gameId,
		Name:  name,
		State: state,
	}

	return game, nil
}
