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

// HasKey checks if a check key is already present
func (p *Proxy) HasKey(key string) bool {
	for i := 0; i < len(p.Checks); i++ {
		if p.Checks[i].Key == key {
			return true
		}
	}
	return false
}
