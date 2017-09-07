package backend

import (
	"fmt"
	"time"
)

// CreateProxy creates a new proxy
func CreateProxy(ip string, port int, t ProxyType, p ProxyProtocol) *Proxy {
	proxy := &Proxy{
		IP:       ip,
		Port:     port,
		Type:     t,
		Protocol: p,
		Joined:   time.Now().Unix(),
	}
	proxy.GenerateIdentifier()
	return proxy
}

// GenerateIdentifier generates the unique identifier for the proxy
func (p *Proxy) GenerateIdentifier() {
	p.Identifier = fmt.Sprintf("%s:%d", p.IP, p.Port)
}
