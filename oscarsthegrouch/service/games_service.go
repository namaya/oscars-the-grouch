package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"namaya/oscarsthegrouch/log"
	"namaya/oscarsthegrouch/model"
	"os"

	"github.com/google/uuid"
)

type GamesService interface {
	ListGames(ctx context.Context) ([]*model.Game, error)
	CreateGame(ctx context.Context, name string) (*model.Game, error)
	ListPlayers(ctx context.Context, gameId string) ([]*model.Player, error)
	GetNominations(ctx context.Context) (*Nominations, error)
	AddPlayer(ctx context.Context, gameId string, user *model.User) (*model.Player, error)
	CreateBallot(ctx context.Context, gameId string, playerId string, votes []*model.Vote) (*model.Ballot, error)
	UpdatePlayerState(ctx context.Context, gameId string, playerId string, state string) (*model.Ballot, error)
	GetMasterBallot(ctx context.Context, gameId string) (*model.Ballot, error)
	UpdateMasterBallot(ctx context.Context, gameId string, votes *model.Vote) (*model.Ballot, error)
	UpdatePlayerScores(ctx context.Context, gameId string, vote *model.Vote) (*model.Ballot, error)
	UpdateGameState(ctx context.Context, gameId string, state string) (*model.Game, error)
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
	logger := log.Get(ctx)

	userId := ctx.Value("userId").(string)
	gameId := uuid.New().String()
	state := "Created"

	// Create game
	result, err := gs.dbClient.ExecContext(ctx, `
		INSERT INTO games (id, name, state, owner_id) VALUES (?, ?, ?, ?)
	`, gameId, name, state, userId)
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

	logger.Debugf("created game '%s'", gameId)

	// Add player to game
	playerId := uuid.New().String()

	result, err = gs.dbClient.ExecContext(ctx, `
		INSERT INTO players (id, game_id, user_id, score, state) VALUES (?, ?, ?, ?, ?)
	`, playerId, gameId, userId, 0, "Waiting")

	if err != nil {
		return nil, fmt.Errorf("CreateGame: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("CreateGame: players: no rows affected")
	}

	logger.Debugf("created player '%s'", playerId)

	// Add master ballot to game
	ballotId := uuid.New().String()
	result, err = gs.dbClient.ExecContext(ctx, `
		INSERT INTO ballots (id, game_id, year) VALUES (?, ?, ?)
	`, ballotId, gameId, "2025")
	if err != nil {
		return nil, fmt.Errorf("CreateGame: %w", err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("CreateGame: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("CreateGame: ballots: no rows affected")
	}

	logger.Debugf("created ballot '%s'", ballotId)

	// NOTE: organizer will add votes to the master ballot during the game

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

type Nominations struct {
	Categories []NominationCategory `json:"categories"`
}

type NominationCategory struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Nominees []Nominee `json:"nominees"`
}

type Nominee struct {
	Work        string `json:"work"`
	Contributor string `json:"contributor"`
}

func (gs *gamesService) GetNominations(ctx context.Context) (*Nominations, error) {
	logger := log.Get(ctx)

	n, err := os.ReadFile("./static/nominations/y2025.json")
	if err != nil {
		return nil, fmt.Errorf("GetNominations: %w", err)
	}

	var nominations Nominations
	err = json.Unmarshal(n, &nominations)
	if err != nil {
		return nil, fmt.Errorf("GetNominations: %w", err)
	}

	logger.Debugf("GetNominations: %v", nominations)

	return &nominations, nil
}

func (gs *gamesService) AddPlayer(ctx context.Context, gameId string, user *model.User) (*model.Player, error) {
	logger := log.Get(ctx)

	playerId := uuid.New().String()

	result, err := gs.dbClient.ExecContext(ctx, `
		INSERT INTO players (id, game_id, user_id, score, state) VALUES (?, ?, ?, ?, ?)
	`, playerId, gameId, user.Id, 0, "Waiting")
	if err != nil {
		return nil, fmt.Errorf("AddPlayer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("AddPlayer: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("AddPlayer: players: no rows affected")
	}

	logger.Debugf("created player '%s'", playerId)

	player := model.Player{
		Id:    playerId,
		User:  user,
		Score: 0,
		State: "Waiting",
	}

	return &player, nil
}

func (gs *gamesService) CreateBallot(ctx context.Context, gameId string, playerId string, votes []*model.Vote) (*model.Ballot, error) {
	logger := log.Get(ctx)
	ballotId := uuid.New().String()
	result, err := gs.dbClient.ExecContext(ctx, `
		INSERT INTO ballots (id, game_id, player_id, year) VALUES (?, ?, ?, ?)
	`, ballotId, gameId, playerId, "2025")

	if err != nil {
		return nil, fmt.Errorf("CreateBallot: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("CreateBallot: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("CreateBallot: ballots: no rows affected")
	}

	ballot := model.Ballot{
		Id:       ballotId,
		PlayerId: playerId,
		Year:     2025,
		Votes:    []*model.Vote{},
	}

	logger.Debugf("created ballot '%s'", ballot)

	// Insert votes
	for _, vote := range votes {
		voteId := uuid.New().String()
		logger.Debugf("creating vote '%s'", voteId)
		result, err = gs.dbClient.ExecContext(ctx, `
			INSERT INTO votes (id, ballot_id, category_id, vote) VALUES (?, ?, ?, ?)
		`, voteId, ballotId, vote.CategoryId, vote.Vote)
		if err != nil {
			return nil, fmt.Errorf("CreateBallot: %w", err)
		}
		rowsAffected, err = result.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("CreateBallot: %w", err)
		}
		if rowsAffected == 0 {
			return nil, fmt.Errorf("CreateBallot: votes: no rows affected")
		}

		ballot.Votes = append(ballot.Votes, &model.Vote{
			CategoryId: vote.CategoryId,
			Vote:       vote.Vote,
		})
	}

	logger.Debugf("added votes")

	return &ballot, nil
}

func (gs *gamesService) UpdatePlayerState(ctx context.Context, gameId string, playerId string, state string) (*model.Ballot, error) {
	logger := log.Get(ctx)

	result, err := gs.dbClient.ExecContext(ctx, `
		UPDATE players SET state = ? WHERE id = ? AND game_id = ?
	`, state, playerId, gameId)
	if err != nil {
		return nil, fmt.Errorf("UpdatePlayerState: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("UpdatePlayerState: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("UpdatePlayerState: players: no rows affected")
	}

	logger.Debugf("updated player '%s' to state '%s'", playerId, state)

	return nil, nil
}

func (gs *gamesService) UpdateMasterBallot(ctx context.Context, gameId string, vote *model.Vote) (*model.Ballot, error) {
	logger := log.Get(ctx)

	voteId := uuid.New().String()

	result, err := gs.dbClient.ExecContext(ctx, `
		INSERT INTO votes (id, ballot_id, category_id, vote)
		SELECT ?, b.id, ?, ?
		FROM ballots b
		WHERE b.game_id = ? AND b.player_id IS NULL;
	`, voteId, vote.CategoryId, vote.Vote, gameId)
	if err != nil {
		return nil, fmt.Errorf("UpdateMasterBallot: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("UpdateMasterBallot: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("UpdateMasterBallot: ballots: no rows affected")
	}

	logger.Debugf("updated master ballot '%s' to vote '%s' for category '%s'", gameId, voteId, vote.CategoryId)

	return nil, nil
}

func (gs *gamesService) GetMasterBallot(ctx context.Context, gameId string) (*model.Ballot, error) {
	logger := log.Get(ctx)

	rows, err := gs.dbClient.QueryContext(ctx, `
		SELECT b.id, b.year, v.category_id, v.vote
		FROM ballots b
			LEFT JOIN votes v ON b.id = v.ballot_id
		WHERE b.game_id = ? AND b.player_id IS NULL;
	`, gameId)
	if err != nil {
		return nil, fmt.Errorf("GetMasterBallot: %w", err)
	}
	defer rows.Close()

	var ballot model.Ballot
	ballot.Id = uuid.New().String()
	ballot.Year = 2025
	ballot.Votes = []*model.Vote{}

	for rows.Next() {
		var vote model.Vote
		var cid sql.NullString
		var v sql.NullInt64

		if err := rows.Scan(&ballot.Id, &ballot.Year, &cid, &v); err != nil {
			return nil, fmt.Errorf("GetMasterBallot: %w", err)
		}

		if cid.Valid {
			vote.CategoryId = cid.String
			vote.Vote = int(v.Int64)
			ballot.Votes = append(ballot.Votes, &vote)
		}
	}

	logger.Debugf("got master ballot '%s'", ballot.Id)

	return &ballot, nil
}

func (gs *gamesService) UpdatePlayerScores(ctx context.Context, gameId string, vote *model.Vote) (*model.Ballot, error) {
	logger := log.Get(ctx)

	result, err := gs.dbClient.ExecContext(ctx, `
		UPDATE players
		SET score = score + 1
		WHERE players.game_id = ? AND players.state = 'Ready' AND players.id IN (
			SELECT b.player_id
			FROM votes v
				RIGHT JOIN ballots b ON v.ballot_id = b.id
			WHERE v.category_id = ? AND v.vote = ?
		);
	`, gameId, vote.CategoryId, vote.Vote)
	if err != nil {
		return nil, fmt.Errorf("UpdatePlayerScores: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("UpdatePlayerScores: %w", err)
	}

	if rowsAffected == 0 {
		logger.Debugf("No correct votes")
		return nil, nil
	}

	logger.Debugf("updated player scores for game '%s'", gameId)

	return nil, nil
}

func (gs *gamesService) UpdateGameState(ctx context.Context, gameId string, state string) (*model.Game, error) {
	logger := log.Get(ctx)

	result, err := gs.dbClient.ExecContext(ctx, `
		UPDATE games SET state = ? WHERE id = ?
	`, state, gameId)
	if err != nil {
		return nil, fmt.Errorf("UpdateGameState: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("UpdateGameState: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("UpdateGameState: games: no rows affected")
	}

	logger.Debugf("updated game '%s' to state '%s'", gameId, state)

	game := &model.Game{
		Id:    gameId,
		State: state,
	}

	return game, nil
}
