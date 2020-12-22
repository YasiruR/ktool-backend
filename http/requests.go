package http

import "github.com/YasiruR/ktool-backend/domain"

//--------------------------add cluster req------------------------------//

type addExistingCluster struct {
	ClusterName  string          `json:"cluster_name"`
	KafkaVersion string          `json:"kafka_version"`
	JmxEnabled   bool            `json:"jmx_enabled"`
	Brokers      []domain.Server `json:"brokers"`
}

type server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

//----------------------------add user req--------------------------------//

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
	SecretName      string `json:"SecretName"`
	UserId          string `json:"UserId"`
	ServiceProvider string `json:"ServiceProvider"`
	Tags            string `json:"Tags"`
	// gke specific
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
	// aws specific
	EksAccessKeyId     string `json:"EksAccessKeyId"`
	EksSecretAccessKey string `json:"EksSecretAccessKey"`
	// azure specific
	AksClientId       string `json:"AksClientId"`
	AksClientSecret   string `json:"AksClientSecret"`
	AksTenantId       string `json:"AksTenantId"`
	AksSubscriptionId string `json:"AksSubscriptionId"`
}

//type SearchSecretsRequest struct {
//	OwnerId         string `json:"OwnerId"`
//	ServiceProvider string `json:"ServiceProvider"`
//}

//-------------------Kubernetes API-----------------------//
