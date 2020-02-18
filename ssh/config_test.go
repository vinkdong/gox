package ssh

import (
	"testing"
)

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
