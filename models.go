package backend

import (
	"github.com/Sirupsen/logrus"
	"github.com/asdine/storm"
	"github.com/labstack/echo"
)

// ProxyProtocol represents the protocol a proxy uses
type ProxyProtocol string

const (
	// HTTPProxyProtocol represents http protocol
	HTTPProxyProtocol = ProxyProtocol("Http")
	// HTTPSProxyProtocol represents https protocol
	HTTPSProxyProtocol = ProxyProtocol("Https")
	// Socks4ProxyProtocol represents socks 4 protocol
	Socks4ProxyProtocol = ProxyProtocol("Socks 4")
	// Socks5ProxyProtocol represents socks 5 protocol
	Socks5ProxyProtocol = ProxyProtocol("Socks 5")
)

// ProxyType represents the anonymity type of the proxy
type ProxyType string

const (
	// HighProxyType is high anonymity
	HighProxyType = ProxyType("High")
	// MiddleProxyType is middle anonymity
	MiddleProxyType = ProxyType("Middle")
	// LowProxyType is low anonymity
	LowProxyType = ProxyType("Low")
)

// Proxy represents a proxy
type Proxy struct {
	ID           int    `storm:"id,increment"`
	Identifier   string `storm:"unique"`
	IP           string
	Port         int
	DeadSince    int
	Alive        bool
	CheckID      int
	Type         ProxyType     `storm:"index"`
	Protocol     ProxyProtocol `storm:"index"`
	Joined       int64
	LastCheck    int64
	Checks       []Check
	ChecksLength int
	SubmittedBy  int
}

// Check represents a proxy check
type Check struct {
	ResponseTime int
	Alive        bool
	DoneBy       int
	Key          string
}

// User represents a user
type User struct {
	ID         int   `storm:"id,increment" json:"-"`
	Submitted  int   `json:"submitted"`
	Checked    int   `json:"checked"`
	Points     int   `json:"points"`
	Joined     int64 `json:"-"`
	LastActive int64 `json:"-"`
}

// Instance represents a server instance
type Instance struct {
	Database       *storm.DB
	Log            *logrus.Logger
	Router         *echo.Echo
	IncomingProxy  chan AddRequest
	IncomingResult chan CheckResult
}

// AddRequest represents a request to add a proxie to the database
type AddRequest struct {
	Proxies []string
	By      int
}

// CheckRequest represents a request for the client
type CheckRequest struct {
	Token    string        `json:"token"`
	IP       string        `json:"ip"`
	Port     int           `json:"port"`
	Protocol ProxyProtocol `json:"protocol"`
}

// CheckResult represents a result from the client
type CheckResult struct {
	Token string `json:"token"`
	Alive bool   `json:"alive"`
	Ms    int    `json:"ms"`
}
