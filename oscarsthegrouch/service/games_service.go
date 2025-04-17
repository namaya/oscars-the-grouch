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
	ListPlayers(ctx context.Context, gameId string) ([]*model.Player, error)
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
		return nil, fmt.Errorf("CreateGame: games: no rows affected")
	}

	result, err = gs.dbClient.ExecContext(ctx, `
		INSERT INTO players (id, game_id, user_id, score, state) VALUES (?, ?, ?, ?, ?)
	`, uuid.New().String(), gameId, userId, 0, "Waiting")

	if err != nil {
		return nil, fmt.Errorf("CreateGame: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("CreateGame: players: no rows affected")
	}

	game := &model.Game{
		Id:    gameId,
		Name:  name,
		State: state,
	}

	return game, nil
}

func (gs *gamesService) ListPlayers(ctx context.Context, gameId string) ([]*model.Player, error) {
	rows, err := gs.dbClient.QueryContext(ctx, `
		SELECT p.id, p.user_id, u.name, u.avatar_uri, p.score, p.state
		FROM players p
			LEFT JOIN users u ON p.user_id = u.id
		WHERE p.game_id = ?
	`, gameId)
	if err != nil {
		return nil, fmt.Errorf("ListPlayers: %w", err)
	}
	defer rows.Close()

	var players []*model.Player
	for rows.Next() {
		var user model.User
		var player model.Player
		if err := rows.Scan(&player.Id, &user.Id, &user.Name, &user.AvatarUri, &player.Score, &player.State); err != nil {
			return nil, fmt.Errorf("ListPlayers: %w", err)
		}
		player.User = &user
		players = append(players, &player)
	}

	return players, nil
}
