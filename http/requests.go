package http

//--------------------------add cluster req------------------------------//

type addExistingCluster struct {
	ClusterName 		string 			`json:"cluster_name"`
	KafkaVersion 		string			`json:"kafka_version"`
	Brokers 			[]server		`json:"brokers"`
}

type server struct {
	Host 		string 			`json:"host"`
	Port 		int				`json:"port"`
}

//----------------------------add user req--------------------------------//

type addUserReq struct {
	Username 		string		`json:"username"`
	Password 		string		`json:"password"`
	AccessLevel 	int			`json:"access_level"`
	FirstName 		string		`json:"first_name"`
	LastName 		string		`json:"last_name"`
	Email 			string 		`json:"email"`
}

type loginUserReq struct {
	Username 		string 		`json:"username"`
	Password 		string		`json:"password"`
}