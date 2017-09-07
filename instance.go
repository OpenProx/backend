package backend

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
)

var ipRegex = regexp.MustCompile("^((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])):(\\d+)$")

// InitAndRun inits and starts the server instance
func (i *Instance) InitAndRun(addr string) error {
	err := i.InitDatabase()
	if err != nil {
		return err
	}

	i.Log = logrus.New()

	i.IncomingProxy = make(chan AddRequest, 100)

	for j := 0; j < 10; j++ {
		go i.IncomingProxyWorker()
	}

	i.InitRouter()

	return i.Router.Start(addr)
}

// IncomingProxyWorker checks incoming proxies
func (i *Instance) IncomingProxyWorker() {
	i.Log.Info("Proxy queue worker started!")
	client := &http.Client{}
	for req := range i.IncomingProxy {
		i.Log.WithFields(logrus.Fields{
			"ID":           req.By,
			"Proxies":      len(req.Proxies),
			"Queue Length": len(i.IncomingProxy),
		}).Info("Starting to process add request...")

		for _, p := range req.Proxies {
			match := ipRegex.FindStringSubmatch(p)
			if len(match) != 7 {
				continue
			}

			ip := match[1]
			port, _ := strconv.Atoi(match[6])
			prot := ProxyProtocol("")

			if i.HasProxy(match[0]) {
				i.Log.WithFields(logrus.Fields{
					"IP":   ip,
					"Port": port,
				}).Warn("Proxy already in database")
				continue
			}

			httpTrans, _ := HTTPTransport(ip, port)
			alive, _, _ := CheckProxy(client, httpTrans, ip, time.Second*6, false)
			if alive {
				prot = HTTPProxyProtocol
			} else {
				httpsTrans, _ := HTTPSTransport(ip, port)
				alive, _, _ = CheckProxy(client, httpsTrans, ip, time.Second*6, true)
				if alive {
					prot = HTTPSProxyProtocol
				} else {
					socks5Trans, _ := Socks5Transport(ip, port)
					alive, _, _ = CheckProxy(client, socks5Trans, ip, time.Second*12, true)
					if alive {
						prot = Socks5ProxyProtocol
					}
				}
			}

			if prot == "" {
				i.Log.WithFields(logrus.Fields{
					"IP":   ip,
					"Port": port,
				}).Warn("Proxy not working...")
				continue
			}

			new := CreateProxy(ip, port, LowProxyType, prot)
			i.Database.Save(new)

			i.Log.WithFields(logrus.Fields{
				"IP":       ip,
				"Port":     port,
				"Protocol": prot,
				"ID":       new.ID,
			}).Info("Proxy working and added!")
		}
	}
}
