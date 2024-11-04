package model

type Ballot struct {
	Categories []Category
	Owner      Player
}

type Category struct {
	Nominees []Nominee
	winnerId string
}

type Nominee struct {
	id   string
	name string
}

type Player struct {
	id   string
	name string
}
