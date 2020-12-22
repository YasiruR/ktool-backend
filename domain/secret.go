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
	ClientId       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	TenantId       string `json:"tenant_id"`
	SubscriptionId string `json:"subscription_id"`
}

type CloudSecret struct {
	UserId          string
	ServiceProvider string
	ID              int
	Name            string
	OwnerId         string
	CreatedOn       string
	CreatedBy       string
	ModifiedOn      string
	ModifiedBy      string
	Activated       bool
	Deleted         bool
	Tags            string
	// gke specific
	GkeType              string
	GkeProjectId         string
	GkePrivateKeyId      string
	GkePrivateKey        string
	GkeClientMail        string
	GkeClientId          string
	GkeAuthUri           string
	GkeTokenUri          string
	GkeAuthX509CertUrl   string
	GkeClientX509CertUrl string
	// aws specific
	EksAccessKeyId     string
	EksSecretAccessKey string
	// azure specific
	AksClientId       string
	AksClientSecret   string
	AksTenantId       string
	AksSubscriptionId string
	//GkeType              string `json:"gke_type"`
	//GkeProjectId         string `json:"gke_project_id"`
	//GkePrivateKeyId      string `json:"gke_private_key_id"`
	//GkePrivateKey        string `json:"gke_private_key"`
	//GkeClientMail        string `json:"gke_client_email"`
	//GkeClientId          string `json:"gke_client_id"`
	//GkeAuthUri           string `json:"gke_auth_uri"`
	//GkeTokenUri          string `json:"gke_token_uri"`
	//GkeAuthX509CertUrl   string `json:"gke_auth_provider_x509_cert_url"`
	//GkeClientX509CertUrl string `json:"gke_client_x509_cert_url"`
	//// aws specific
	//EksAccessKeyId     string `json:"eks_access_key_id"`
	//EksSecretAccessKey string `json:"eks_secret_access_key"`
	//// azure specific
	//AksClientId       string `json:"aks_client_id"`
	//AksClientSecret   string `json:"aks_client_secret"`
	//AksTenantId       string `json:"aks_tenant_id"`
	//AksSubscriptionId string `json:"aks_subscription_id"`
	Validate bool
}

type Result struct {
	SecretList []Secret
	Status     int
	Message    string
	Error      error
	ErrorMsg   string
}

type DAOResult struct {
	Secret  CloudSecret
	Status  int
	Message string
	Error   error
}

type Validation struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

//type SecretDAO interface {
//	AddSecret(ctx context.Context, addSecretRequest *http.AddSecretRequest) (result Result)
//	DeleteSecret(ctx context.Context, secretId string) (result Result)
//	GetAllSecretsByUser(ctx context.Context, userId string) (result Result)
//}
