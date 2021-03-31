package dns

import (
	"context"
	"net"
	"strings"
	"time"
)

type DNS struct {
	Server   string
	Resolver *net.Resolver
}

func New(server string) (*DNS, error) {
	// 添加端口
	if !strings.Contains(server, ":") && server != "" {
		server = server + ":53"
	}

	dial := func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: time.Millisecond * time.Duration(10000),
		}
		return d.DialContext(ctx, "udp", server)
	}
	if server == "" {
		dial = nil
	}
	return &DNS{Server: server,
		Resolver: &net.Resolver{
			Dial: dial,
		}}, nil
}

func (dns *DNS) LookUpNS(domain string) []string {
	ns, _ := dns.Resolver.LookupNS(context.Background(), domain)
	nsList := make([]string, len(ns), len(ns))
	for i, n := range ns {
		nsList[i] = n.Host
	}
	return nsList
}

func (dns *DNS) LookupIPAddr(domain string) []string {
	ns, _ := dns.Resolver.LookupIPAddr(context.Background(), domain)
	ipList := make([]string, len(ns), len(ns))
	for i, n := range ns {
		ipList[i] = n.IP.String()
	}
	return ipList
}
