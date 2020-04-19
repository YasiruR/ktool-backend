package http

type addExistingCluster struct {
	ClusterName  string   `json:"cluster_name"`
	KafkaVersion string   `json:"kafka_version"`
	Brokers      []server `json:"brokers"`
}

type server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
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
	Username    string `json:"username"`
	Password    string `json:"password"`
	AccessLevel int    `json:"access_level"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
}

type loginUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//-------------------------secret-management-api------------------//
type addSecretRequest struct {
	secretName string
	userId     string
	service    string
	keyType    int
	key        string
	tags       string
}
