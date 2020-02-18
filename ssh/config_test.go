package ssh

import (
	"testing"
)

var c = Config{
	HostConfig{
		Host:                  "abc1",
		ServerAliveInterval:   60,
		HostName:              "abc1.com",
		User:                  "root1",
		Port:                  22,
		StrictHostKeyChecking: "",
		UserKnownHostsFile:    "",
		IdentityFile:          "",
		LogLevel:              "",
	},
	HostConfig{
		Host:                  "abc2",
		ServerAliveInterval:   160,
		HostName:              "abc2.com",
		User:                  "root2",
		Port:                  22,
		StrictHostKeyChecking: "",
		UserKnownHostsFile:    "",
		IdentityFile:          "",
		LogLevel:              "",
	},
	HostConfig{
		Host:                  "abc3",
		ServerAliveInterval:   180,
		HostName:              "abc3.com",
		User:                  "root3",
		Port:                  0,
		StrictHostKeyChecking: "",
		UserKnownHostsFile:    "",
		IdentityFile:          "",
		LogLevel:              "",
	},
}

func TestConfig_Marshal(t *testing.T) {
	data := c.Marshal()
	cfg := &Config{}
	Unmarshal(cfg, data)
	if len(*cfg) != 3 {
		t.Fatal("Marshal config failed")
	}
}

func TestUnmarshal(t *testing.T) {
	data := `Host     *

ServerAliveInterval 60 

Host    abc
HostName vinkdong.com
Port 22
User root
Host efg
HostName gobuildrun.com
Port 22
User root`

	c := &Config{}
	Unmarshal(c, []byte(data))
	if len(*c) != 3 {
		t.Fatal("Unmarshal config failed")
	}
}

func TestConfig_SetHost(t *testing.T) {

	newHost := HostConfig{
		Host:                  "abc4",
		ServerAliveInterval:   200,
		HostName:              "abcUser",
		User:                  "",
		Port:                  2222,
		StrictHostKeyChecking: "",
		UserKnownHostsFile:    "",
		IdentityFile:          "",
		LogLevel:              "",
	}
	c.SetHost(newHost)
	if len(c) != 4 {
		t.Fatal("set config host failed")
	}
	newHost = HostConfig{
		Host:                  "abc3",
		ServerAliveInterval:   200,
		HostName:              "abc3User",
		User:                  "",
		Port:                  5555,
		StrictHostKeyChecking: "",
		UserKnownHostsFile:    "",
		IdentityFile:          "",
		LogLevel:              "",
	}
	c.SetHost(newHost)
	if c.GetHost("abc3").Port != 5555 {
		t.Fatal("set config host failed, should get port 5555 of abc3 host")
	}
}
