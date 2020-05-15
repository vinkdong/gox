package net

import (
	"fmt"
	"github.com/vinkdong/gox/log"
	"net"
	"strings"
)

/*
return bindPort
*/
func SafeListen(host string, port int32, acceptFunc func(conn net.Conn)) (conn net.Listener, bindPort int32, err error) {
	bindPort = port
	// 防止非root用户无权限绑定port
	if port < 8000 {
		port += 8000
	}
start:
	in := fmt.Sprintf("%s:%d", host, port)
	c, err := net.Listen("tcp", in)
	if err != nil {

		if strings.Contains(err.Error(), "address already in use") {
			port++
			log.Warnf("port %d in use,change to %d", port, port+1)
			goto start
		}
		return c, port, err
	}
	ch := make(chan error, 0)
	go func() {
		for {
			c, err := c.Accept()
			if err != nil {
				ch <- err
			}
			go acceptFunc(c)
		}
	}()
	return c, port, nil
}
