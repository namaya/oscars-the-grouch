package model

type User struct {
	Id        string
	Name      string
	AvatarUri string
}

type Player struct {
	Id    string
	User  *User
	Score int
	State string
}

type Game struct {
	Id      string
	Name    string
	State   string
	OwnerId string
	// Players []Player
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
