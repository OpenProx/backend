package backend

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestHTTPProxy(t *testing.T) {
	ip := "212.237.34.237"
	trans, _ := HTTPTransport(ip, 3128)
	alive, ms, err := CheckProxy(http.DefaultClient, trans, ip, time.Second*6, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(alive, ms, err)
}

func TestHTTPSProxy(t *testing.T) {
	ip := "212.237.34.237"
	trans, _ := HTTPSTransport(ip, 3128)
	alive, ms, err := CheckProxy(http.DefaultClient, trans, ip, time.Second*6, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(alive, ms, err)
}

func TestSocks5Proxy(t *testing.T) {
	ip := "216.138.34.140"
	trans, _ := Socks5Transport(ip, 38834)
	alive, ms, err := CheckProxy(http.DefaultClient, trans, ip, time.Second*6, false)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(alive, ms, err)
}
