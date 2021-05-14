package internal

import "testing"

func TestRemoteAddrs(t *testing.T) {
	cases := map[string]string{
		"[::1]:1234":      "127.0.0.1",       // Golang default :1234 syntax uses ipv6
		"localhost:1234":  "127.0.0.1",       // Always assume loopback and bypass /etc/hosts
		"127.0.0.1:1234":  "127.0.0.1",       // Always assume loopback
		"1.2.3.4:1234":    "1.2.3.4",         // Assume host without port
		"1.2.3.4:":        "1.2.3.4",         // Assume host without delimiter
		"1.2.3.4":         "1.2.3.4",         // Assume host
		"www.nivenly.com": "www.nivenly.com", // Assume host
	}
	for given, expected := range cases {
		actual := RemoteAddrToHost(given)
		if actual != expected {
			t.Errorf("Expected (%s) Actual (%s)", expected, actual)
		}
	}
}
