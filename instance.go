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
	i.IncomingResult = make(chan CheckResult, 1000)

	for j := 0; j < 10; j++ {
		go i.IncomingProxyWorker()
	}
	go i.IncomingResultWorker()

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

// IncomingResultWorker checks incoming results
func (i *Instance) IncomingResultWorker() {
	for chk := range i.IncomingResult {
		pid, uid, chkid, err := DecodeRequestToken(chk.Token)
		if err != nil {
			i.Log.WithError(err).Warn("Error while decoding result token...")
			continue
		}

		// Get Proxy from db
		var proxy Proxy
		if err := i.Database.One("ID", pid, &proxy); err != nil {
			i.Log.WithFields(logrus.Fields{
				"ProxyID": pid,
			}).WithError(err).Warn("Proxy not in database...")
			continue
		}

		// Check if CheckID matches
		if proxy.CheckID != chkid {
			i.Log.WithFields(logrus.Fields{
				"Got":     chkid,
				"Current": proxy.CheckID,
			}).Warn("Proxy check too late...")
			continue
		}

		// Check if key is already used
		if proxy.HasUserCheck(uid) {
			i.Log.WithFields(logrus.Fields{
				"ID": uid,
			}).Warn("Proxy check duplication...")
			continue
		}

		// Finalize
		proxy.Checks = append(proxy.Checks, Check{chk.Ms, chk.Alive, uid})
		proxy.ChecksLength++

		if proxy.ChecksLength >= 5 {
			// Count Alive
			a := 0
			for i := 0; i < len(proxy.Checks); i++ {
				if proxy.Checks[i].Alive {
					a++
				}
			}

			// Alive
			if a >= 3 {
				proxy.Alive = true
				proxy.DeadSince = 0
				// Reward
			} else {
				proxy.Alive = false
				proxy.DeadSince++
			}

			// Next Check
			proxy.LastCheck = time.Now().Unix()
			proxy.CheckID++
			proxy.Checks = make([]Check, 0)

			i.Log.WithFields(logrus.Fields{
				"Proxy":  proxy.Identifier,
				"Result": proxy.Alive,
			}).Info("Proxy check finished!")
		} else {
			i.Log.WithFields(logrus.Fields{
				"Proxy":  proxy.Identifier,
				"Result": chk.Alive,
			}).Info("Proxy check result recieved")
		}

		i.Database.Update(&proxy)
	}
}
