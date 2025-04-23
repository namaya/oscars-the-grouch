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
}

type Ballot struct {
	Id       string
	PlayerId string
	Year     int
	Votes    []*Vote
}

type Vote struct {
	CategoryId string
	Vote       int
}

type Category struct {
	Id       string
	Name     string
	Nominees []Nominee
}

type Nominee struct {
	Id           string
	Work         string
	Contributors string
}
