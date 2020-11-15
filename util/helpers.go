package util

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/YasiruR/ktool-backend/domain"
)

func StrToVMType(strVmType string) (vmType containerservice.VMSizeTypes) {
	return containerservice.StandardA1
}

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

func StringPointerListToEscapedCSV(list *[]*string) (csv string) {
	l := *list
	if len(l) > 0 {
		csv = ""
		for i := 0; i < len(*list); i++ {
			val := l[i]
			csv = fmt.Sprintf(csv+"%q, ", *val)
		}
		return csv[:len(csv)-2]
	}
	return ""
}
