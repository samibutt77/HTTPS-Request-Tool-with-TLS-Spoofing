package httpclient
import "golang.org/x/net/http2"


import (
	"os"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
	"encoding/base64"

	utls "github.com/refraction-networking/utls"
	goproxy "golang.org/x/net/proxy"
	"abtls/internal/tlsprofile"
)

// 🔐 Compute JA3 hash from TLS fingerprint
func computeJA3Hash(spec utls.ClientHelloSpec) string {
	var b strings.Builder
	b.WriteString("771") // TLS 1.2

	// Cipher Suites
	b.WriteString(",")
	for i, cs := range spec.CipherSuites {
		b.WriteString(fmt.Sprintf("%d", cs))
		if i < len(spec.CipherSuites)-1 {
			b.WriteString("-")
		}
	}

	// Extensions
	b.WriteString(",")
	var extIDs []int
	for _, ext := range spec.Extensions {
		switch ext.(type) {
		case *utls.SNIExtension:
			extIDs = append(extIDs, 0)
		case *utls.StatusRequestExtension:
			extIDs = append(extIDs, 5)
		case *utls.SupportedCurvesExtension:
			extIDs = append(extIDs, 10)
		case *utls.SupportedPointsExtension:
			extIDs = append(extIDs, 11)
		case *utls.SignatureAlgorithmsExtension:
			extIDs = append(extIDs, 13)
		case *utls.ALPNExtension:
			extIDs = append(extIDs, 16)
		case *utls.SCTExtension:
			extIDs = append(extIDs, 18)
		case *utls.SessionTicketExtension:
			extIDs = append(extIDs, 35)
		case *utls.PSKKeyExchangeModesExtension:
			extIDs = append(extIDs, 45)
		case *utls.KeyShareExtension:
			extIDs = append(extIDs, 51)
		case *utls.SupportedVersionsExtension:
			extIDs = append(extIDs, 43)
		case *utls.CookieExtension:
			extIDs = append(extIDs, 44)
		case *utls.RenegotiationInfoExtension:
			extIDs = append(extIDs, 65281)
		}
	}
	for i, id := range extIDs {
		b.WriteString(fmt.Sprintf("%d", id))
		if i < len(extIDs)-1 {
			b.WriteString("-")
		}
	}

	// Elliptic Curves
	b.WriteString(",")
	var curves []int
	for _, ext := range spec.Extensions {
		if e, ok := ext.(*utls.SupportedCurvesExtension); ok {
			for _, c := range e.Curves {
				curves = append(curves, int(c))
			}
		}
	}
	for i, c := range curves {
		b.WriteString(fmt.Sprintf("%d", c))
		if i < len(curves)-1 {
			b.WriteString("-")
		}
	}

	// EC Point Formats
	b.WriteString(",")
	var points []int
	for _, ext := range spec.Extensions {
		if e, ok := ext.(*utls.SupportedPointsExtension); ok {
			for _, p := range e.SupportedPoints {
				points = append(points, int(p))
			}
		}
	}
	for i, p := range points {
		b.WriteString(fmt.Sprintf("%d", p))
		if i < len(points)-1 {
			b.WriteString("-")
		}
	}

	// Final MD5 hash
	sum := md5.Sum([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

// 🧠 Check if JA3 hash exists in known_ja3.txt
func IsKnownJA3(ja3 string) bool {
	data, err := os.ReadFile("known_JA3.txt")
	if err != nil {
		return false
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == ja3 {
			return true
		}
	}
	return false
}

func SaveJA3Hash(ja3 string) {
    f, err := os.OpenFile("known_JA3.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("❌ Could not open known_JA3.txt:", err)
        return
    }
    defer f.Close()

    // Add newline if file doesn't already end with one
    stat, _ := f.Stat()
    if stat.Size() > 0 {
        f.WriteString("\n")
    }

    if _, err := f.WriteString(ja3 + "\n"); err != nil {
        fmt.Println("❌ Failed to write JA3:", err)
    }
}

// 🧠 Builds HTTP client with proxy, uTLS and cookie jar
func NewClient(proxyStr, proxyType, tlsProfile, username, password string) (*http.Client, error) {
	var transport *http.Transport

	// Function to wrap a uTLS connection
	dialTLS := func(rawConn net.Conn, addr string, profile string) (net.Conn, error) {
		host := addr
		if strings.Contains(addr, ":") {
			host, _, _ = net.SplitHostPort(addr)
		}
		spec, _ := tlsprofile.GetClientHello(profile)
		utlsConn := utls.UClient(rawConn, &utls.Config{ServerName: host}, utls.HelloCustom)
		if err := utlsConn.ApplyPreset(&spec); err != nil {
			return nil, err
		}
		if err := utlsConn.Handshake(); err != nil {
			return nil, err
		}

		hash := computeJA3Hash(spec)
		fmt.Println("🔑 JA3 Hash:", hash)

		return utlsConn, nil
	}

	if proxyType == "socks5" {
		// SOCKS5 Dialer with optional auth
		var auth *goproxy.Auth
		if username != "" && password != "" {
			auth = &goproxy.Auth{
				User:     username,
				Password: password,
			}
		}
		socksDialer, err := goproxy.SOCKS5("tcp", proxyStr, auth, &net.Dialer{Timeout: 45 * time.Second})
		if err != nil {
			return nil, err
		}

		transport = &http.Transport{
			DialTLS: func(network, addr string) (net.Conn, error) {
				rawConn, err := socksDialer.Dial(network, addr)
				if err != nil {
					return nil, err
				}
				return dialTLS(rawConn, addr, tlsProfile)
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

	} else {
		// HTTP proxy with optional basic auth
		proxyURL, err := url.Parse("http://" + proxyStr)
		if err != nil {
			return nil, err
		}

		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialTLS: func(network, addr string) (net.Conn, error) {
				rawConn, err := net.DialTimeout(network, addr, 45*time.Second)
				if err != nil {
					return nil, err
				}
				return dialTLS(rawConn, addr, tlsProfile)
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		// Add basic auth header for HTTP proxies
		if username != "" && password != "" {
			auth := username + ":" + password
			basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
			transport.ProxyConnectHeader = make(http.Header)
			transport.ProxyConnectHeader.Set("Proxy-Authorization", basic)
		}
	}

	// Enable HTTP/2 if ALPN negotiated
	_ = http2.ConfigureTransport(transport)
	
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Transport: transport,
		Timeout:   45 * time.Second,
		Jar:       jar,
	}
	return client, nil
}



func GetJA3HashForProfile(profile string) string {
    spec, _ := tlsprofile.GetClientHello(profile)
    return computeJA3Hash(spec)
}

