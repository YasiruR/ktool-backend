package kafka

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"github.com/YasiruR/ktool-backend/log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func generateRSAKeys(ctx context.Context, clusterName string) (caFile, certFile, keyFile string, err error) {
	//creating a certificate authority
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"K-Tool"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,	//implies that this is our ca
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	//generating rsa keys
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		log.Logger.TraceContext(ctx, "failed to generate rsa keys", err)
		return
	}

	publicKey := key.PublicKey

	//generate the certificate authority bytes
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &key.PublicKey, key)
	if err != nil {
		return
	}

	//create certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"K-Tool"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},	//we want this certificate to be valid in localhost
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	//generate certbytes
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &key.PublicKey, key)
	if err != nil {
		return
	}

	//------------saving keys, cert and ca-----------------//

	absPath, err := filepath.Abs("../src/")
	if err != nil {
		log.Logger.ErrorContext(ctx, "converting relative path to absolute path failed", err)
		return
	}

	err = os.Mkdir(absPath + "/" + clusterName, os.ModePerm)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("creating key folder failed for cluster - %v", clusterName), err)
		return
	}

	absPath = absPath + "/" + clusterName + "/"

	//save ca
	err = saveCA(ctx, absPath + clusterName + "_ca.pem", caBytes)
	if err != nil {
		return
	}

	err = saveCert(ctx, absPath + clusterName + "_cert.pem", certBytes)
	if err != nil {
		return
	}

	//save private key
	err = saveGobKey(ctx, absPath + clusterName + "_private.key", key)
	if err != nil {
		return
	}
	err = savePEMKey(ctx, absPath + clusterName + "_private.pem", key)
	if err != nil {
		return
	}

	//save public key
	err = saveGobKey(ctx, absPath + clusterName + "_public.key", publicKey)
	if err != nil {
		return
	}
	err = savePublicPEMKey(ctx, absPath + clusterName + "_public.pem", publicKey)
	if err != nil {
		return
	}

	log.Logger.TraceContext(ctx, "rsa key pair, certificate and certificate authority are generated and saved for the cluster", clusterName)

	return absPath + clusterName + "_ca.pem", absPath + clusterName + "_cert.pem", absPath + clusterName + "_private.pem", nil
}


func saveCA(ctx context.Context, fileName string, caBytes []byte) (err error) {
	outFile, err := os.Create(fileName)
	if err != nil {
		log.Logger.ErrorContext(ctx, "creating file for rsa keys failed", err)
		return err
	}
	//caPEM := new(bytes.Buffer)

	err = pem.Encode(outFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding ca failed", err)
		return err
	}

	return nil
}

func saveCert(ctx context.Context, fileName string, certBytes []byte) (err error) {
	outFile, err := os.Create(fileName)
	if err != nil {
		log.Logger.ErrorContext(ctx, "creating file for rsa keys failed", err)
		return err
	}

	err = pem.Encode(outFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding cert failed", err)
		return err
	}

	return nil
}

func savePEMKey(ctx context.Context, fileName string, key *rsa.PrivateKey) (err error) {
	outFile, err := os.Create(fileName)
	if err != nil {
		log.Logger.ErrorContext(ctx, "creating file for rsa keys failed", err)
		return err
	}

	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding rsa private key failed", err)
		return err
	}

	return nil
}

func savePublicPEMKey(ctx context.Context, fileName string, pubkey rsa.PublicKey) (err error) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	if err != nil {
		log.Logger.ErrorContext(ctx, "marshalling rsa key failed", err)
		return err
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	if err != nil {
		log.Logger.ErrorContext(ctx, "creating file for rsa public key failed", err)
	}
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding rsa public key failed", err)
		return err
	}

	return nil
}

func saveGobKey(ctx context.Context, fileName string, key interface{}) (err error) {
	outFile, err := os.Create(fileName)
	if err != nil {
		log.Logger.ErrorContext(ctx, "creating file for rsa key failed", err)
		return err
	}
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	if err != nil {
		log.Logger.ErrorContext(ctx, "encoding rsa key failed", err)
		return err
	}

	return nil
}