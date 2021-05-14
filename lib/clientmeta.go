package lib

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Ullaakut/nmap/v2"

	"github.com/kris-nova/logger"
)

type ClientMeta struct {
	RemoteAddress   string
	LastPortScan    time.Time
	PublicPortProbe map[string]string
	PortScan        *nmap.Run
}

var cache = map[string]*ClientMeta{}

func GetClientMeta(r *http.Request) *ClientMeta {
	// When coming from a reverse proxy we pass the
	// original host IP via a proxy header "X-Real-IP"
	// so we use that for our reverse port probe.
	clientAddr := r.Header.Get("X-Real-IP")
	logger.Info("Client: %s", clientAddr)
	//for name, values := range r.Header {
		//logger.Info(name)
		//logger.Info(strings.Join(values, ","))
		//logger.Info("----")
	//}

	clientMeta := &ClientMeta{
		RemoteAddress: clientAddr,
	}
	if cached, ok := cache[clientAddr]; ok {
		logger.Debug("Cached client meta for: %s", clientAddr)
		clientMeta = cached
	}


	// Each client meta has a scanner
	go func() {
		scanner, err := nmap.NewScanner(
			nmap.WithTargets(RemoteAddrToHost(clientAddr)),
			nmap.WithContext(context.TODO()),
		)
		if err != nil {
			logger.Warning("Port scanning error: %v", err)
			return
		}
		results, warnings, err := scanner.Run()
		if err != nil {
			logger.Warning("Port scanning error: %v", err)
			return
		}
		if warnings != nil {
			logger.Warning("Port scanning warning: %v", err)
			return
		}
		clientMeta.PortScan = results
		logger.Debug("Client (%s) hosts: %d", clientAddr, len(results.Hosts))
		for _, host := range results.Hosts {
			logger.Debug("Host (%s) ports: %d", host.Comment, len(host.Ports))
			for _, port := range host.Ports {
				clientMeta.PublicPortProbe[fmt.Sprintf("%d %s", port.ID, port.Protocol)] = fmt.Sprintf("[%s] %s", port.State, port.Service)
			}
		}
		clientMeta.LastPortScan = time.Now()
	}()
	cache[clientAddr] = clientMeta
	return clientMeta
}

func RemoteAddrToHost(remoteAddr string) string {
	if strings.Contains(remoteAddr, "::1") {
		// Return the LO interface
		return "127.0.0.1"
	}
	if strings.Contains(remoteAddr, "localhost") {
		// Return the LO interface
		return "127.0.0.1"
	}
	if strings.Contains(remoteAddr, "127.0.0.1") {
		// Return the LO interface
		return "127.0.0.1"
	}
	if strings.Contains(remoteAddr, ":") {
		spl := strings.Split(remoteAddr, ":")
		if len(spl) <= 2 {
			// 1.2.3.4:80
			// 1.2.3.4:
			// return 1.2.3.4
			return spl[0]
		}
	}
	return remoteAddr
}
