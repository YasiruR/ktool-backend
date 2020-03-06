package http

type addExistingCluster struct {
	ClusterName 		string 			`json:"cluster_name"`
	KafkaVersion 		string			`json:"kafka_version"`
	Brokers 			[]server		`json:"brokers"`
}

type server struct {
	Host 		string 			`json:"host"`
	Port 		int				`json:"port"`
}

//type testCluster struct {
//	Brokers 	[]server 		`json:"brokers"`
//}
//type connectToCluster struct {
//	ClusterID 		int				`json:"cluster_id"`
//	ClusterName 	string			`json:"cluster_name"`
//	Brokers 		[]string		`json:"brokers"`
//}

//----------------------------addUserReq--------------------------------//

type addUserReq struct {
	Username 		string		`json:"username"`
	Password 		string		`json:"password"`
	AccessLevel 	int			`json:"access_level"`
}

type loginUserReq struct {
	Username 		string 		`json:"username"`
	Password 		string		`json:"password"`
}