package proxy

import (
	"bufio"
	"os"
	"strings"
)

// Proxy holds a proxy address and its type (http, socks5, or mixed)
type Proxy struct {
	Address  string
	Type     string // always "mixed" in this case
	Username string
	Password string
	Full     string 
}

// LoadProxyList loads proxies from a file and marks them as "mixed"
func LoadProxyList(path string) ([]Proxy, error) {
	var proxies []Proxy

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var proxy Proxy
		proxy.Type = "mixed"
		proxy.Full = line

		if strings.Contains(line, "@") {
			// Handle user:pass@ip:port
			parts := strings.SplitN(line, "@", 2)
			auth := strings.SplitN(parts[0], ":", 2)
			if len(auth) == 2 {
				proxy.Username = auth[0]
				proxy.Password = auth[1]
				proxy.Address = parts[1]
			} else {
				continue // malformed auth
			}
		} else {
			// No authentication
			proxy.Address = line
		}

		proxies = append(proxies, proxy)
	}

	return proxies, nil
}

// ✅ Add support for calling it with a mixed-type loader
func LoadMixedProxyList(path string) ([]Proxy, error) {
	return LoadProxyList(path) // assumes all are "mixed"
}
