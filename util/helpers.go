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

func StrToVMType(str string) (vmType containerservice.VMSizeTypes) {

	if str == "StandardA1V2" {
		return containerservice.StandardA1V2
	}
	if str == "StandardA2V2" {
		return containerservice.StandardA2V2
	}
	if str == "StandardA4V2" {
		return containerservice.StandardA4V2
	}
	if str == "StandardA8V2" {
		return containerservice.StandardA8V2
	}
	if str == "StandardA2mV2" {
		return containerservice.StandardA2mV2
	}
	if str == "StandardA4mV2" {
		return containerservice.StandardA4mV2
	}
	if str == "StandardA8mV2" {
		return containerservice.StandardA8mV2
	}
	//if str == "StandardDc1sV2" {
	//	return containerservice.StandardDc1sV2
	//}
	//if str == "StandardDc2sV2" {
	//	return containerservice.StandardDc2sV2
	//}
	//if str == "StandardDc4sV2" {
	//	return containerservice.StandardDc4sV2
	//}
	//if str == "StandardDc8V2" {
	//	return containerservice.StandardDc8V2
	//}
	if str == "StandardD1V2" {
		return containerservice.StandardD1V2
	}
	if str == "StandardD2V2" {
		return containerservice.StandardD2V2
	}
	if str == "StandardD3V2" {
		return containerservice.StandardD3V2
	}
	if str == "StandardD4V2" {
		return containerservice.StandardD4V2
	}
	if str == "StandardD5V2" {
		return containerservice.StandardD5V2
	}
	//if str == "StandardDs1V2" {
	//	return containerservice.StandardDs1V2
	//}
	//if str == "StandardDs2V2" {
	//	return containerservice.StandardDs2V2
	//}
	//if str == "StandardDs3V2" {
	//	return containerservice.StandardDs3V2
	//}
	//if str == "StandardDs4V2" {
	//	return containerservice.StandardDs4V2
	//}
	//if str == "StandardDs5V2" {
	//	return containerservice.StandardDs5V2
	//}
	if str == "StandardD2V3" {
		return containerservice.StandardD2V3
	}
	if str == "StandardD4V3" {
		return containerservice.StandardD4V3
	}
	if str == "StandardD8V3" {
		return containerservice.StandardD8V3
	}
	if str == "StandardD16V3" {
		return containerservice.StandardD16V3
	}
	if str == "StandardD32V3" {
		return containerservice.StandardD32V3
	}
	//if str == "StandardD48V3" {
	//	return containerservice.StandardD48V3
	//}
	if str == "StandardD64V3" {
		return containerservice.StandardD64V3
	}
	if str == "StandardD2sV3" {
		return containerservice.StandardD2sV3
	}
	if str == "StandardD4sV3" {
		return containerservice.StandardD4sV3
	}
	if str == "StandardD8sV3" {
		return containerservice.StandardD8sV3
	}
	if str == "StandardD16sV3" {
		return containerservice.StandardD16sV3
	}
	if str == "StandardD32sV3" {
		return containerservice.StandardD32sV3
	}
	//if str == "StandardD48sV3" {
	//	return containerservice.StandardD48sV3
	//}
	if str == "StandardD64sV3" {
		return containerservice.StandardD64sV3
	}
	//if str == "StandardD2aV4" {
	//	return containerservice.StandardD2aV4
	//}
	//if str == "StandardD4aV4" {
	//	return containerservice.StandardD4aV4
	//}
	//if str == "StandardD8aV4" {
	//	return containerservice.StandardD8aV4
	//}
	//if str == "StandardD16aV4" {
	//	return containerservice.StandardD16aV4
	//}
	//if str == "StandardD32aV4" {
	//	return containerservice.StandardD32aV4
	//}
	//if str == "StandardD48aV4" {
	//	return containerservice.StandardD48aV4
	//}
	//if str == "StandardD64aV4" {
	//	return containerservice.StandardD64aV4
	//}
	//if str == "StandardD96aV4" {
	//	return containerservice.StandardD96aV4
	//}
	//if str == "StandardD2asV4" {
	//	return containerservice.StandardD2asV4
	//}
	//if str == "StandardD4asV4" {
	//	return containerservice.StandardD4asV4
	//}
	//if str == "StandardD8asV4" {
	//	return containerservice.StandardD8asV4
	//}
	//if str == "StandardD16asV4" {
	//	return containerservice.StandardD16asV4
	//}
	//if str == "StandardD32asV4" {
	//	return containerservice.StandardD32asV4
	//}
	//if str == "StandardD48asV4" {
	//	return containerservice.StandardD48asV4
	//}
	//if str == "StandardD64asV4" {
	//	return containerservice.StandardD64asV4
	//}
	//if str == "StandardD96asV4" {
	//	return containerservice.StandardD96asV4
	//}
	//if str == "StandardD2dV4" {
	//	return containerservice.StandardD2dV4
	//}
	//if str == "StandardD4dV4" {
	//	return containerservice.StandardD4dV4
	//}
	//if str == "StandardD8dV4" {
	//	return containerservice.StandardD8dV4
	//}
	//if str == "StandardD16dV4" {
	//	return containerservice.StandardD16dV4
	//}
	//if str == "StandardD32dV4" {
	//	return containerservice.StandardD32dV4
	//}
	//if str == "StandardD48dV4" {
	//	return containerservice.StandardD48dV4
	//}
	//if str == "StandardD64dV4" {
	//	return containerservice.StandardD64dV4
	//}
	//if str == "StandardD2dsV4" {
	//	return containerservice.StandardD2dsV4
	//}
	//if str == "StandardD4dsV4" {
	//	return containerservice.StandardD4dsV4
	//}
	//if str == "StandardD8dsV4" {
	//	return containerservice.StandardD8dsV4
	//}
	//if str == "StandardD16dsV4" {
	//	return containerservice.StandardD16dsV4
	//}
	//if str == "StandardD32dsV4" {
	//	return containerservice.StandardD32dsV4
	//}
	//if str == "StandardD48dsV4" {
	//	return containerservice.StandardD48dsV4
	//}
	//if str == "StandardD64dsV4" {
	//	return containerservice.StandardD64dsV4
	//}
	//if str == "StandardD2V4" {
	//	return containerservice.StandardD2V4
	//}
	//if str == "StandardD4V4" {
	//	return containerservice.StandardD4V4
	//}
	//if str == "StandardD8V4" {
	//	return containerservice.StandardD8V4
	//}
	//if str == "StandardD16V4" {
	//	return containerservice.StandardD16V4
	//}
	//if str == "StandardD32V4" {
	//	return containerservice.StandardD32V4
	//}
	//if str == "StandardD48V4" {
	//	return containerservice.StandardD48V4
	//}
	//if str == "StandardD64V4" {
	//	return containerservice.StandardD64V4
	//}
	//if str == "StandardD2sV4" {
	//	return containerservice.StandardD2sV4
	//}
	//if str == "StandardD4sV4" {
	//	return containerservice.StandardD4sV4
	//}
	//if str == "StandardD8sV4" {
	//	return containerservice.StandardD8sV4
	//}
	//if str == "StandardD16sV4" {
	//	return containerservice.StandardD16sV4
	//}
	//if str == "StandardD32sV4" {
	//	return containerservice.StandardD32sV4
	//}
	//if str == "StandardD48sV4" {
	//	return containerservice.StandardD48sV4
	//}
	//if str == "StandardD64sV4" {
	//	return containerservice.StandardD64sV4
	//}
	if str == "StandardF2sV2" {
		return containerservice.StandardF2sV2
	}
	if str == "StandardF4sV2" {
		return containerservice.StandardF4sV2
	}
	if str == "StandardF8sV2" {
		return containerservice.StandardF8sV2
	}
	if str == "StandardF16sV2" {
		return containerservice.StandardF16sV2
	}
	if str == "StandardF32sV2" {
		return containerservice.StandardF32sV2
	}
	//if str == "StandardF48sV2" {
	//	return containerservice.StandardF48sV2
	//}
	if str == "StandardF64sV2" {
		return containerservice.StandardF64sV2
	}
	if str == "StandardD11V2" {
		return containerservice.StandardD11V2
	}
	if str == "StandardD12V2" {
		return containerservice.StandardD12V2
	}
	if str == "StandardD13V2" {
		return containerservice.StandardD13V2
	}
	if str == "StandardD14V2" {
		return containerservice.StandardD14V2
	}
	if str == "StandardE2V3" {
		return containerservice.StandardE2V3
	}
	if str == "StandardE4V3" {
		return containerservice.StandardE4V3
	}
	if str == "StandardE8V3" {
		return containerservice.StandardE8V3
	}
	if str == "StandardE16V3" {
		return containerservice.StandardE16V3
	}
	//if str == "StandardE20V3" {
	//	return containerservice.StandardE20V3
	//}
	if str == "StandardE32V3" {
		return containerservice.StandardE32V3
	}
	//if str == "StandardE48V3" {
	//	return containerservice.StandardE48V3
	//}
	if str == "StandardE64V3" {
		return containerservice.StandardE64V3
	}
	if str == "StandardE2sV3" {
		return containerservice.StandardE2sV3
	}
	//if str == "StandardE20sV3" {
	//	return containerservice.StandardE20sV3
	//}
	//if str == "StandardE2aV4" {
	//	return containerservice.StandardE2aV4
	//}
	//if str == "StandardE4aV4" {
	//	return containerservice.StandardE4aV4
	//}
	//if str == "StandardE8aV4" {
	//	return containerservice.StandardE8aV4
	//}
	//if str == "StandardE16aV4" {
	//	return containerservice.StandardE16aV4
	//}
	//if str == "StandardE20aV4" {
	//	return containerservice.StandardE20aV4
	//}
	//if str == "StandardE32aV4" {
	//	return containerservice.StandardE32aV4
	//}
	//if str == "StandardE48aV4" {
	//	return containerservice.StandardE48aV4
	//}
	//if str == "StandardE64aV4" {
	//	return containerservice.StandardE64aV4
	//}
	//if str == "StandardE96aV4" {
	//	return containerservice.StandardE96aV4
	//}
	//if str == "StandardE2asV4" {
	//	return containerservice.StandardE2asV4
	//}
	//if str == "StandardE4asV4" {
	//	return containerservice.StandardE4asV4
	//}
	//if str == "StandardE8asV4" {
	//	return containerservice.StandardE8asV4
	//}
	//if str == "StandardE16asV4" {
	//	return containerservice.StandardE16asV4
	//}
	//if str == "StandardE20asV4" {
	//	return containerservice.StandardE20asV4
	//}
	//if str == "StandardE32asV4" {
	//	return containerservice.StandardE32asV4
	//}
	//if str == "StandardE48asV4" {
	//	return containerservice.StandardE48asV4
	//}
	//if str == "StandardE64asV4" {
	//	return containerservice.StandardE64asV4
	//}
	//if str == "StandardE2dV4" {
	//	return containerservice.StandardE2dV4
	//}
	//if str == "StandardE4dV4" {
	//	return containerservice.StandardE4dV4
	//}
	//if str == "StandardE8dV4" {
	//	return containerservice.StandardE8dV4
	//}
	//if str == "StandardE16dV4" {
	//	return containerservice.StandardE16dV4
	//}
	//if str == "StandardE20dV4" {
	//	return containerservice.StandardE20dV4
	//}
	//if str == "StandardE32dV4" {
	//	return containerservice.StandardE32dV4
	//}
	//if str == "StandardE48dV4" {
	//	return containerservice.StandardE48dV4
	//}
	//if str == "StandardE64dV4" {
	//	return containerservice.StandardE64dV4
	//}
	//if str == "StandardE2dsV4" {
	//	return containerservice.StandardE2dsV4
	//}
	//if str == "StandardE4dsV4" {
	//	return containerservice.StandardE4dsV4
	//}
	//if str == "StandardE8dsV4" {
	//	return containerservice.StandardE8dsV4
	//}
	//if str == "StandardE16dsV4" {
	//	return containerservice.StandardE16dsV4
	//}
	//if str == "StandardE20dsV4" {
	//	return containerservice.StandardE20dsV4
	//}
	//if str == "StandardE32dsV4" {
	//	return containerservice.StandardE32dsV4
	//}
	//if str == "StandardE48dsV4" {
	//	return containerservice.StandardE48dsV4
	//}
	//if str == "StandardE2V4" {
	//	return containerservice.StandardE2V4
	//}
	//if str == "StandardE4V4" {
	//	return containerservice.StandardE4V4
	//}
	//if str == "StandardE8V4" {
	//	return containerservice.StandardE8V4
	//}
	//if str == "StandardE16V4" {
	//	return containerservice.StandardE16V4
	//}
	//if str == "StandardE32V4" {
	//	return containerservice.StandardE32V4
	//}
	//if str == "StandardE48V4" {
	//	return containerservice.StandardE48V4
	//}
	//if str == "StandardE64V4" {
	//	return containerservice.StandardE64V4
	//}
	//if str == "StandardE2sV4" {
	//	return containerservice.StandardE2sV4
	//}
	//if str == "StandardE4sV4" {
	//	return containerservice.StandardE4sV4
	//}
	//if str == "StandardE8sV4" {
	//	return containerservice.StandardE8sV4
	//}
	//if str == "StandardE16sV4" {
	//	return containerservice.StandardE16sV4
	//}
	//if str == "StandardE20sV4" {
	//	return containerservice.StandardE20sV4
	//}
	//if str == "StandardE32sV4" {
	//	return containerservice.StandardE32sV4
	//}
	//if str == "StandardE48sV4" {
	//	return containerservice.StandardE48sV4
	//}
	//if str == "StandardM8ms" {
	//	return containerservice.StandardM8ms
	//}
	//if str == "StandardM16ms" {
	//	return containerservice.StandardM16ms
	//}
	//if str == "StandardM32ts" {
	//	return containerservice.StandardM32ts
	//}
	//if str == "StandardM32ls" {
	//	return containerservice.StandardM32ls
	//}
	//if str == "StandardM32ms" {
	//	return containerservice.StandardM32ms
	//}
	//if str == "StandardL8sV2" {
	//	return containerservice.StandardL8sV2
	//}
	//if str == "StandardL16sV2" {
	//	return containerservice.StandardL16sV2
	//}
	//if str == "StandardL32sV2" {
	//	return containerservice.StandardL32sV2
	//}
	//if str == "StandardL48sV2" {
	//	return containerservice.StandardL48sV2
	//}
	//if str == "StandardL64sV2" {
	//	return containerservice.StandardL64sV2
	//}
	//if str == "StandardNc6" {
	//	return containerservice.StandardNc6
	//}
	//if str == "StandardNc12" {
	//	return containerservice.StandardNc12
	//}
	//if str == "StandardNc24" {
	//	return containerservice.StandardNc24
	//}
	//if str == "StandardNc6sV2" {
	//	return containerservice.StandardNc6sV2
	//}
	//if str == "StandardNc12sV2" {
	//	return containerservice.StandardNc12sV2
	//}
	//if str == "StandardNc24sV2" {
	//	return containerservice.StandardNc24sV2
	//}
	//if str == "StandardNc6sV3" {
	//	return containerservice.StandardNc6sV3
	//}
	//if str == "StandardNc12sV3" {
	//	return containerservice.StandardNc12sV3
	//}
	//if str == "StandardNc24sV3" {
	//	return containerservice.StandardNc24sV3
	//}
	//if str == "StandardNc4asT4V3" {
	//	return containerservice.StandardNc4asT4V3
	//}
	//if str == "StandardNc8asT4V3" {
	//	return containerservice.StandardNc8asT4V3
	//}
	//if str == "StandardNc16asT4V3" {
	//	return containerservice.StandardNc16asT4V3
	//}
	//if str == "StandardNc64asT4V3" {
	//	return containerservice.StandardNc64asT4V3
	//}
	//if str == "StandardNd6s" {
	//	return containerservice.StandardNd6s
	//}
	//if str == "StandardNd12s" {
	//	return containerservice.StandardNd12s
	//}
	//if str == "StandardNd24s" {
	//	return containerservice.StandardNd24s
	//}
	//if str == "StandardNd40rsV2" {
	//	return containerservice.StandardNd40rsV2
	//}
	//if str == "StandardNv6" {
	//	return containerservice.StandardNv6
	//}
	//if str == "StandardNv12" {
	//	return containerservice.StandardNv12
	//}
	//if str == "StandardNv24" {
	//	return containerservice.StandardNv24
	//}
	//if str == "StandardNv12sV3" {
	//	return containerservice.StandardNv12sV3
	//}
	//if str == "StandardNv24sV3" {
	//	return containerservice.StandardNv24sV3
	//}
	//if str == "StandardNv48sV3" {
	//	return containerservice.StandardNv48sV3
	//}
	if str == "StandardH8" {
		return containerservice.StandardH8
	}
	if str == "StandardH16" {
		return containerservice.StandardH16
	}
	if str == "StandardH8m" {
		return containerservice.StandardH8m
	}
	if str == "StandardH16m" {
		return containerservice.StandardH16m
	}
	//if str == "StandardHb60rs" {
	//	return containerservice.StandardHb60rs
	//}
	//if str == "StandardHb120rsV2" {
	//	return containerservice.StandardHb120rsV2
	//}
	//if str == "StandardHc44rs" {
	//	return containerservice.StandardHc44rs
	//}
	return containerservice.StandardF2sV2 // set default machine type here
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

//func EncodePrivateKey(privateKey *rsa.PrivateKey) (string, error) {
//	pubKeyBytes := ssh.MarshalAuthorizedKey(privateKey)
//
//	log.Println("Private key encoded to string")
//	return string(pubKeyBytes), nil
//}

// writePemToFile writes keys to a file
func WriteKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}
