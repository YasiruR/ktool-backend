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
type AddSecretRequest struct {
	SecretName       string `json:"SecretName"`
	UserId           string `json:"UserId"`
	ServiceProvider  string `json:"ServiceProvider"`
	Tags             string `json:"Tags"`
	GkeType          string `json:"GkeType"`
	GkeProjectId     string `json:"GkeProjectId"`
	GkePrivateKeyId  string `json:"GkePrivateKeyId"`
	GkePrivateKey    string `json:"GkePrivateKey"`
	GkeClientMail    string `json:"GkeClientMail"`
	GkeClientId      string `json:"GkeClientId"`
	GkeAuthUri       string `json:"GkeAuthUri"`
	GkeTokenUri      string `json:"GkeTokenUri"`
	GkeAuthCertUrl   string `json:"GkeAuthCertUrl"`
	GkeClientCertUrl string `json:"GkeClientCertUrl"`
}

type SearchSecretsRequest struct {
	OwnerId         string `json:"OwnerId"`
	ServiceProvider string `json:"ServiceProvider"`
}
