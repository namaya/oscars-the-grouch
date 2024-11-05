package service

import (
	"database/sql"
	"namaya/oscarsthegrouch/model"
)

type BallotsService interface {
}

type ballotsService struct {
}

func NewBallotsService() BallotsService {
	return &ballotsService{}
}

func (bs *ballotsService) SaveBallot(b *model.Ballot) error {
	db, err := sql.Open("sqlite3", "./oscarsthegrouch.db")
	defer db.Close()

	if err != nil {
		return err
	}

	db.Exec("INSERT INTO ballots (name, email, movie) VALUES (?, ?)", b.Owner, b.Categories)

	return nil
}
