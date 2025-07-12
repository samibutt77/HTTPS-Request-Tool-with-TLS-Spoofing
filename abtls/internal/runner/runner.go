package runner

import (
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/andybalholm/brotli"

	"abtls/internal/config"
	"abtls/internal/httpclient"
	"abtls/internal/proxy"
)

func Run(cfg *config.Config) {
	proxies, err := proxy.LoadMixedProxyList("proxies.txt")
	if err != nil || len(proxies) == 0 {
		fmt.Println("❌ Error loading proxies or empty list.")
		return
	}

	presets := []string{"chrome", "firefox", "safari"}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(proxies), func(i, j int) {
		proxies[i], proxies[j] = proxies[j], proxies[i]
	})

	for idx, p := range proxies {
		profile := cfg.TLSProfile
		if profile == "random" {
			profile = presets[rand.Intn(len(presets))]
		}

		var client *http.Client
		var err error
		var proxyTypeUsed string

		// Try HTTP first
		client, err = httpclient.NewClient(p.Address, "http", profile, p.Username, p.Password)
		proxyTypeUsed = "http"

		fmt.Printf("\n[%d/%d] Trying proxy (%s): %s | TLS profile: %s\n", idx+1, len(proxies), proxyTypeUsed, p.Address, profile)

		if err != nil {
			fmt.Printf("\n[%d/%d] ❌ Failed to build HTTP client: %s\n", idx+1, len(proxies), err)
		} else {
			req, _ := http.NewRequest("GET", cfg.URL, nil)
			req.Header = httpclient.GetHeadersForProfile(profile)

			fmt.Println("Request Headers:")
			for key, values := range req.Header {
				for _, v := range values {
					fmt.Printf("%s: %s\n", key, v)
				}
			}

			resp, err := client.Do(req)
			if err == nil {
				if handleResponse(resp, req, client, profile) {
					return
				}
			} else {
				fmt.Println("❌ Request failed on HTTP:", err)
			}
		}

		// Fallback to SOCKS5
		client, err = httpclient.NewClient(p.Address, "socks5", profile, p.Username, p.Password)
		proxyTypeUsed = "socks5"

		fmt.Printf("\n[%d/%d] Trying proxy (%s): %s | TLS profile: %s\n", idx+1, len(proxies), proxyTypeUsed, p.Address, profile)

		if err != nil {
			fmt.Printf("\n[%d/%d] ❌ Failed to build SOCKS5 client: %s\n", idx+1, len(proxies), err)
			continue
		}

		req, _ := http.NewRequest("GET", cfg.URL, nil)
		req.Header = httpclient.GetHeadersForProfile(profile)

		fmt.Println("Request Headers:")
		for key, values := range req.Header {
			for _, v := range values {
				fmt.Printf("%s: %s\n", key, v)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("❌ Request failed on SOCKS5:", err)
		} else {
			if handleResponse(resp, req, client, profile) {
				return
			}
		}

		delay := time.Duration(rand.Intn(cfg.MaxDelay-cfg.MinDelay)+cfg.MinDelay) * time.Millisecond
		fmt.Printf("⏱ Waiting %v before next proxy...\n", delay)
		time.Sleep(delay)
	}

	fmt.Println("❌ All proxies failed or none returned 200 OK without challenge.")
}

func handleResponse(resp *http.Response, req *http.Request, client *http.Client, profile string) bool {
	defer resp.Body.Close()

	var reader io.ReadCloser
	var err error
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("❌ Failed to decode gzip body:", err)
			return false
		}
		defer reader.Close()
	case "br":
		reader = io.NopCloser(brotli.NewReader(resp.Body))
	default:
		reader = resp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("❌ Error reading response body:", err)
		return false
	}

	snippet := string(body)
	if len(snippet) > 200 {
		snippet = snippet[:200]
	}

	challengeMarkers := []string{"cf-challenge", "g-recaptcha", "verify you are human"}
	challenged := false
	for _, marker := range challengeMarkers {
		if strings.Contains(strings.ToLower(snippet), marker) {
			challenged = true
			fmt.Printf("⚠️ Challenge detected in response body (marker: %s)\n", marker)
			break
		}
	}

	if resp.StatusCode == 200 && !challenged {
		fmt.Println("✅ Success")
		fmt.Println("Status:", resp.StatusCode)
		fmt.Println("Body Snippet:", snippet)

		hash := httpclient.GetJA3HashForProfile(profile)
		fmt.Println("🔑 JA3 Hash (from profile):", hash)

		if !httpclient.IsKnownJA3(hash) {
			httpclient.SaveJA3Hash(hash)
			fmt.Println("✅ JA3 added to known_JA3.txt")
		}

		if client.Jar != nil {
			cookies := client.Jar.Cookies(req.URL)
			if len(cookies) > 0 {
				fmt.Println("🍪 Cookies after response:")
				for _, c := range cookies {
					fmt.Printf("  %s = %s\n", c.Name, c.Value)
				}
			}
		}

		return true
	} else if resp.StatusCode == 403 {
		fmt.Println("🚫 Blocked with status 403")
	} else if challenged {
		fmt.Println("🚧 Blocked due to challenge in body")
	} else if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		fmt.Printf("🔁 Detected redirect to: %s\n", resp.Header.Get("Location"))
	}
	return false
}
