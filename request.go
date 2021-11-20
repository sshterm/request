package request

import (
	"context"
	"github.com/ncruces/go-dns"
	"net"
	"net/http"
	"time"
)

var (
	dots = map[string][]string{
		"dns.alidns.com": {"223.5.5.5", "223.6.6.6", "2400:3200::1", "2400:3200:baba::1"},
		"dot.pub":        {"162.14.21.178", "162.14.21.56"},
	}
)

func NewRequest() *http.Client {
	client := &http.Client{}
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err == nil {
			ips, err := GetIp(host)
			if err == nil && len(ips) > 0 {
				ip := ips[0].String()
				addr = net.JoinHostPort(ip, port)
			}
		}
		return dialer.DialContext(ctx, network, addr)
	}
	return client
}

func GetIp(host string) (ips []net.IPAddr, err error) {
	for server, addresses := range dots {
		var resolver *net.Resolver
		resolver, err = dns.NewDoTResolver(server, dns.DoTAddresses(addresses...), dns.DoTCache())
		if err == nil {
			ips, err = resolver.LookupIPAddr(context.TODO(), host)
			if err == nil {
				break
			}
		}
	}
	return
}
