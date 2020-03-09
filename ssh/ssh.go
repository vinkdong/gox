package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"time"
)

type SSH struct {
	Host       string `json:"host"`
	Username   string `json:"username"`
	Port       int    `json:"port"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	Stdout     io.Writer
	Stderr     io.Writer
	Stdin      io.Reader
}

func (s *SSH) newSession() (*ssh.Session, error) {
	// todo close client
	client, err := s.NewClient()
	if err != nil {
		return nil, err
	}
	return client.NewSession()
}

func (s *SSH) NewClient() (*ssh.Client, error) {
	// todo: auth by public keys
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(s.Password))

	clientConfig := &ssh.ClientConfig{
		User:    s.Username,
		Auth:    auth,
		Timeout: time.Minute * 2,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	return ssh.Dial("tcp", addr, clientConfig)
}
