package backend

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// BaseChecker represents the url of the ip checking api
const BaseChecker = "ipv4bot.whatismyipaddress.com"

// CheckProxy checks if a proxy is working
func CheckProxy(client *http.Client, trans *http.Transport, ip string, timeout time.Duration, ssl bool) (bool, int64, error) {
	client.Transport = trans
	client.Timeout = timeout
	url := "http://" + BaseChecker
	if ssl {
		url = "https://" + BaseChecker
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, -1, err
	}
	now := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(now)
	if err != nil {
		return false, elapsed.Nanoseconds() / 1000000, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, elapsed.Nanoseconds() / 1000000, err
	}
	defer resp.Body.Close()
	return strings.Trim(string(body), " -\n\r\t") == ip, elapsed.Nanoseconds() / 1000000, nil
}

// HTTPTransport creates a http proxy transport layer
func HTTPTransport(ip string, port int) (*http.Transport, error) {
	url, err := url.Parse(fmt.Sprintf("http://%s:%d", ip, port))
	if err != nil {
		return nil, err
	}
	return &http.Transport{Proxy: http.ProxyURL(url)}, nil
}

// HTTPSTransport creates a https proxy transport layer
func HTTPSTransport(ip string, port int) (*http.Transport, error) {
	url, err := url.Parse(fmt.Sprintf("https://%s:%d", ip, port))
	if err != nil {
		return nil, err
	}
	return &http.Transport{Proxy: http.ProxyURL(url), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, nil
}

// Socks5Transport creates a socks 5 proxy transport layer
func Socks5Transport(ip string, port int) (*http.Transport, error) {
	dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", ip, port), nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	return &http.Transport{Dial: dialer.Dial, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, nil
}
