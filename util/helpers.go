package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice"
	"github.com/YasiruR/ktool-backend/domain"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
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

//crypto utils
//func main() {
//	savePrivateFileTo := "./id_rsa_test"
//	savePublicFileTo := "./id_rsa_test.pub"
//	bitSize := 4096
//
//	privateKey, err := generatePrivateKey(bitSize)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//
//	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//
//	privateKeyBytes := encodePrivateKeyToPEM(privateKey)
//
//	err = writeKeyToFile(privateKeyBytes, savePrivateFileTo)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//
//	err = writeKeyToFile([]byte(publicKeyBytes), savePublicFileTo)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//}

// generatePrivateKey creates a RSA Private Key of specified byte size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func GeneratePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

// writePemToFile writes keys to a file
func WriteKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}
