package cloud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/service"
	"github.com/sparrc/go-ping"
	"golang.org/x/crypto/ssh"
)

func ConnectToServer(ctx context.Context, ipAddress string) (err error) {
	config := &ssh.ClientConfig{
		User:	"username",
		Auth: 	[]ssh.AuthMethod{ssh.PublicKeys(Key)},
	}

	addr := ipAddress + ":22"

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Logger.ErrorContext(ctx, "dialing tcp connection to server failed", err)
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		log.Logger.ErrorContext(ctx, "creating server session failed", err)
		return err
	}

	SessionList = append(SessionList, session)

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run("pwd"); err != nil {
		log.Logger.ErrorContext(ctx, "command failed", err)
		return err
	}

	return nil
}

func pingToServer(ctx context.Context, ipAddress string) (ok bool, err error) {
	pinger, err := ping.NewPinger(ipAddress)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("ping to server %s", ipAddress), err)
		return false, err
	}

	fmt.Println("about to")

	//pinger.SetPrivileged(true)

	pinger.Count = service.Cfg.PingRetry
	pinger.Run()
	stats := pinger.Statistics()
	fmt.Println("stats : ", stats)

	if stats.PacketsSent == 0 {
		//log.Logger.ErrorContext(ctx, "packets sent : 0")
		return false, errors.New("could not send packets to ping")
	}

	if stats.PacketsRecv != stats.PacketsSent {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("packets sent : %v, packets received : %v", stats.PacketsSent, stats.PacketsRecv))
		return false, errors.New("could not receive all the packets")
	}

	return true, nil
}