package domain

type Secret struct {
	ID         int
	Name       string
	OwnerId    string
	Provider   string
	Type       string
	CreatedOn  string
	CreatedBy  string
	ModifiedOn string
	ModifiedBy string
	Activated  bool
	Deleted    bool
	Tags       string
}

type GkeSecret struct {
	Type              string `json:"type"`
	ProjectId         string `json:"project_id"`
	PrivateKeyId      string `json:"private_key_id"`
	PrivateKey        string `json:"private_key"`
	ClientMail        string `json:"client_email"`
	ClientId          string `json:"client_id"`
	AuthUri           string `json:"auth_uri"`
	TokenUri          string `json:"token_uri"`
	AuthX509CertUrl   string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl string `json:"client_x509_cert_url"`
}

type EksSecret struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
	Token  string `json:"token"`
	Region string `json:"region"`
}

type AksSecret struct {
}

type CloudSecret struct {
	ID         int
	Name       string
	OwnerId    string
	Provider   string
	CreatedOn  string
	CreatedBy  string
	ModifiedOn string
	ModifiedBy string
	Activated  bool
	Deleted    bool
	Tags       string
	// gke specific
	Type              string `json:"type"`
	ProjectId         string `json:"project_id"`
	PrivateKeyId      string `json:"private_key_id"`
	PrivateKey        string `json:"private_key"`
	ClientMail        string `json:"client_email"`
	ClientId          string `json:"client_id"`
	AuthUri           string `json:"auth_uri"`
	TokenUri          string `json:"token_uri"`
	AuthX509CertUrl   string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl string `json:"client_x509_cert_url"`
	// aws specific
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	// azure specific
}

type Result struct {
	SecretList []Secret
	Status     int
	Message    string
	Error      error
}

type DAOResult struct {
	SecretList []CloudSecret
	//SecretList []GkeSecret
	Status  int
	Message string
	Error   error
}

//type SecretDAO interface {
//	AddSecret(ctx context.Context, addSecretRequest *http.AddSecretRequest) (result Result)
//	DeleteSecret(ctx context.Context, secretId string) (result Result)
//	GetAllSecretsByUser(ctx context.Context, userId string) (result Result)
//}
