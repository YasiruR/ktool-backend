package util

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/ktool-backend/domain"
)

func ConvertSecretToGKESecretBytes(secret domain.CloudSecret) (gkeSecret []byte, err error) {
	return json.Marshal(domain.GkeSecret{
		Type:              secret.GkeType,
		ProjectId:         secret.GkeProjectId,
		PrivateKeyId:      secret.GkePrivateKeyId,
		PrivateKey:        secret.GkePrivateKey,
		ClientMail:        secret.GkeClientMail,
		ClientId:          secret.GkeClientId,
		AuthUri:           secret.GkeAuthUri,
		TokenUri:          secret.GkeTokenUri,
		AuthX509CertUrl:   secret.GkeAuthX509CertUrl,
		ClientX509CertUrl: secret.GkeClientX509CertUrl,
	})
}

func StringListToEscapedCSV(list []string) (csv string) {
	if len(list) > 0 {
		csv = ""
		for i := 0; i < len(list); i++ {
			csv = fmt.Sprintf(csv+"%q, ", list[i])
		}
		return csv[:len(csv)-2]
	}
	return ""
}
