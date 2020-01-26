package cloud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/YasiruR/ktool-backend/log"
	"github.com/YasiruR/ktool-backend/service"
	"github.com/reiver/go-telnet"
	"github.com/sparrc/go-ping"
	"golang.org/x/crypto/ssh"
	"time"
)

func PingToServer(ctx context.Context, ipAddress string) (ok bool, err error) {

	var stats *ping.Statistics
	pinger, err := ping.NewPinger(ipAddress)
	if err != nil {
		log.Logger.ErrorContext(ctx, fmt.Sprintf("ping to server %s", ipAddress), err)
		return false, err
	}

	//pinger.SetPrivileged(true)

	pinged := make(chan bool)

	go func() {
		pinger.Count = service.Cfg.PingRetry
		pinger.Run()
		stats = pinger.Statistics()

		if stats.PacketsSent == 0 {
			log.Logger.ErrorContext(ctx, "could not send any packet")
			pinged <- false
		}

		if stats.PacketsRecv != stats.PacketsSent {
			log.Logger.ErrorContext(ctx, fmt.Sprintf("packets sent : %v, packets received : %v", stats.PacketsSent, stats.PacketsRecv))
			pinged <- false
		}

		pinged <- true
	}()

	select {
	case res := <- pinged:
		if res {
			log.Logger.TraceContext(ctx, fmt.Sprintf("ping successful : %v", ipAddress))
			return true, nil
		} else {
			return false, errors.New("server ping failed")
		}
	case <- time.After(time.Duration(int64(service.Cfg.PingTimeout))*time.Second):
		log.Logger.ErrorContext(ctx, fmt.Sprintf("%v server ping timeout : %v seconds", ipAddress, service.Cfg.PingTimeout))
		return false, errors.New("server ping timeout")
	}
}

func TelnetToPort(ctx context.Context, ipAddress string, port int) (ok bool, err error) {

	done := make(chan bool)

	address := fmt.Sprintf("%s:%v", ipAddress, port)
	go func() {
		_, err = telnet.DialTo(address)
		if err != nil {
			log.Logger.ErrorContext(ctx, fmt.Sprintf("telnet to server %s and port %v failed", ipAddress, port), err)
			done <- false
		}
		done <- true
	}()

	select {
	case res := <- done:
		if res {
			log.Logger.TraceContext(ctx, fmt.Sprintf("telnet successful : %s", address))
			return true, nil
		} else {
			return false, errors.New("telnet failed")
		}
	case <- time.After(time.Duration(int64(service.Cfg.PingTimeout))*time.Second):
		log.Logger.ErrorContext(ctx, fmt.Sprintf("telnet timeout : %v, address : %v", service.Cfg.PingTimeout, address))
		return false, errors.New("telnet failed")
	}
}

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
