package service

import (
	"context"
	"database/sql"
	"namaya/oscarsthegrouch/model"
)

type BallotsService interface {
	SaveBallot(ctx context.Context, ballot *model.Ballot) error
}

type ballotsService struct {
	dbClient *sql.DB
}

func NewBallotsService(dbClient *sql.DB) BallotsService {
	return &ballotsService{
		dbClient: dbClient,
	}
}

func (bs *ballotsService) SaveBallot(ctx context.Context, ballot *model.Ballot) error {
	// bs.dbClient.Exec("INSERT INTO ballots (name, email, movie) VALUES (?, ?)", ballot.Owner, ballot.Categories)
	return nil
}
