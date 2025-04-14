package model

type User struct {
	Id   string
	Name string
	// Avatar string
}

type Player struct {
	Id   string
	Name string
}

type Game struct {
	Id      string
	Name    string
	State   string
	OwnerId User
	Players []Player
}

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
