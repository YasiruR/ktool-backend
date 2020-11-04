package cloud

import (
	"github.com/YasiruR/ktool-backend/log"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os/user"
)

var (
	Key				ssh.Signer
	SessionList		[]*ssh.Session
)

func Init() {
	key, err := getSSHKey()
	if err != nil {
		log.Logger.Error("error occurred while initializing the ssh connection", err)
	}
	Key = key
}

func getSSHKey() (key ssh.Signer, err error) {
	usr, err := user.Current()
	if err != nil {
		log.Logger.Fatal("failed in fetching the user for ssh retrieve", err)
		return nil, err
	}

	file := usr.HomeDir + "/.ssh/ktool"
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Logger.Fatal("failed in reading the ssh file", err)
		return nil, err
	}

	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		log.Logger.Fatal("error occurred in parsing private key", err)
		return nil, err
	}

	log.Logger.Trace("ssh key fetched from the user")

	return key, nil
}
