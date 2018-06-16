package block

import (
	"net"
	"sync"
	"time"

	"github.com/pkg/sftp"

	"github.com/CardInfoLink/log"
	"golang.org/x/crypto/ssh"
)

var finishChan = make(chan int, 1)

func blockWrap(f func(*ssh.Client, int), timeout int, sleep int) (isSuccess bool) {
	conn := sshDial(timeout)
	if conn == nil {
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Debug(r)
			}
		}()
		wg.Done()

		ticker := time.NewTicker(time.Duration(timeout) * time.Second)
		select {
		case <-ticker.C:
			log.Debug("receive timeout ticker")
			conn.Close() // close connection, so that using `conn` in block() will report error or panic
			isSuccess = false
		case <-finishChan:
			isSuccess = true
		}
	}()

	f(conn, sleep)

	wg.Wait()
	return
}

func block(conn *ssh.Client, sleep int) {
	log.Debug("block begin")
	time.Sleep(time.Duration(sleep) * time.Second) // block for sleep seconds

	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Errorf("get sftp client error: %v", err)
		return
	}
	defer client.Close()

	client.Open("/tmp")

	finishChan <- 1
	log.Debug("block finished")
}

func sshDial(timeout int) *ssh.Client {
	var authMethods []ssh.AuthMethod
	authMethods = append(authMethods, ssh.Password("xxx"))
	config := &ssh.ClientConfig{
		User:    "webapp",
		Auth:    authMethods,
		Timeout: time.Duration(timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	conn, err := ssh.Dial("tcp", "xxx:22", config)
	if err != nil {
		log.Errorf("connect to sftp error: %v", err)
		return nil
	}
	return conn
}
