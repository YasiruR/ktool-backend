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
	Type              string
	ProjectId         string
	PrivateKeyId      string
	PrivateKey        string
	ClientMail        string
	ClientId          string
	AuthUri           string
	TokenUri          string
	AuthX509CertUrl   string
	ClientX509CertUrl string
}

type EksSecret struct {
	AccessKeyId     string
	SecretAccessKey string
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
	Type              string
	ProjectId         string
	PrivateKeyId      string
	PrivateKey        string
	ClientMail        string
	ClientId          string
	AuthUri           string
	TokenUri          string
	AuthX509CertUrl   string
	ClientX509CertUrl string
	// aws specific
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
	Status     int
	Message    string
	Error      error
}

//type SecretDAO interface {
//	AddSecret(ctx context.Context, addSecretRequest *http.AddSecretRequest) (result Result)
//	DeleteSecret(ctx context.Context, secretId string) (result Result)
//	GetAllSecretsByUser(ctx context.Context, userId string) (result Result)
//}
