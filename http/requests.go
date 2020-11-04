package http

import "github.com/YasiruR/ktool-backend/domain"

//--------------------------add cluster req------------------------------//

type addExistingCluster struct {
	ClusterName 		string 			`json:"cluster_name"`
	KafkaVersion 		string			`json:"kafka_version"`
	JmxEnabled 			bool			`json:"jmx_enabled"`
	Brokers 			[]domain.Server	`json:"brokers"`
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