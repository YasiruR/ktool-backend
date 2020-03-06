package domain

var LoggedInUsers []User

type User struct {
	Id 				int
	Username		string
	Token			string
	AccessLevel 	int
}
