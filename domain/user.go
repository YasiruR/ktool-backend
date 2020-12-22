package domain

var LoggedInUserMap map[int]User

type User struct {
	Id 					int
	Username			string
	Token				string
	AccessLevel 		int
	ConnectedClusters	[]KCluster
}
