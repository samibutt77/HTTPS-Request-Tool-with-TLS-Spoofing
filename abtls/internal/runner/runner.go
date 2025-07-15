package runner

import (
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"os"

	"github.com/andybalholm/brotli"

	"abtls/internal/config"
	"abtls/internal/httpclient"
	"abtls/internal/proxy"
)


func saveSuccessfulProxy(proxyStr, profile string) {
	f, err := os.OpenFile("successful_combos.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("❌ Could not write successful combos:", err)
		return
	}
	defer f.Close()

	entry := fmt.Sprintf("%s | profile: %s\n", proxyStr, profile)
	f.WriteString(entry)
}


func Run(cfg *config.Config) {
	proxies, err := proxy.LoadMixedProxyList("proxies.txt")
	if err != nil || len(proxies) == 0 {
		fmt.Println("❌ Error loading proxies or empty list.")
		return
	}

	presets := []string{"chrome", "firefox", "safari"}
	//rand.Seed(time.Now().UnixNano())
	if cfg.Shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(proxies), func(i, j int) {
			proxies[i], proxies[j] = proxies[j], proxies[i]
		})
	}


	var targetRequests int
	fmt.Print("🔢 Enter the total number of requests to send: ")
	fmt.Scanf("%d", &targetRequests)

	if targetRequests <= 0 {
		fmt.Println("❌ Invalid number of requests. Exiting.")
		return
	}

	successCount := 0
	totalSent := 0

	writtenCombos := make(map[string]bool) // ✅ deduplication

	fmt.Printf("\n🚀 Starting with target of %d requests...\n", targetRequests)

	for idx, p := range proxies {
		if totalSent >= targetRequests {
			break
		}

		profile := cfg.TLSProfile
		if profile == "random" {
			profile = presets[rand.Intn(len(presets))]
		}

		client, err := httpclient.NewClient(p.Address, "http", profile, p.Username, p.Password)
		if err != nil {
			fmt.Printf("❌ Failed to build HTTP client for proxy %s: %s\n", p.Address, err)
			continue
		}

		fmt.Printf("\n[%d/%d] 🌐 Using proxy: %s@%s | TLS profile: %s\n",
			idx+1, len(proxies), p.Username, p.Address, profile)

		comboKey := fmt.Sprintf("%s|%s", p.Full, profile)
		wroteThisCombo := false

		// Keep using this proxy until it fails or max requests reached
		for totalSent < targetRequests {
			req, _ := http.NewRequest("GET", cfg.URL, nil)
			headers := httpclient.GetOrderedHeadersForProfile(profile)
			for _, h := range headers {
				req.Header.Add(h.Key, h.Value)
			}

			fmt.Println("Request Headers:")
			for _, h := range headers {
				fmt.Printf("%s: %s\n", h.Key, h.Value)
			}

			resp, err := client.Do(req)
			totalSent++

			if err != nil {
				fmt.Println("❌ Request failed:", err)
				break
			}

			if handleResponse(resp, req, client, profile, p.Full) {
				successCount++

				// ✅ Save only first time
				if !wroteThisCombo && !writtenCombos[comboKey] {
					saveSuccessfulProxy(p.Full, profile)
					writtenCombos[comboKey] = true
					wroteThisCombo = true
				}
			} else {
				break // Stop using this proxy if response isn't 200 OK or has challenge
			}

			delay := time.Duration(rand.Intn(cfg.MaxDelay-cfg.MinDelay)+cfg.MinDelay) * time.Millisecond
			fmt.Printf("⏱ Waiting %v before next request on same proxy...\n", delay)
			time.Sleep(delay)
		}
	}

	fmt.Printf("\n✅ Finished. Total Requests: %d | Successes: %d | Success Rate: %.2f%%\n",
		totalSent, successCount, (float64(successCount)/float64(totalSent))*100)
}





func handleResponse(resp *http.Response, req *http.Request, client *http.Client, profile string, fullProxy string) bool {
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

	// 🚧 Heuristic Challenge Detection (Cloudflare / Akamai / PerimeterX etc.)
	challengeMarkers := []string{
		"cf-challenge",
		"g-recaptcha",
		"verify you are human",
		"location.reload",                     // Reload on challenge solve
		"XMLHttpRequest.prototype.send",      // JS override seen in Akamai
		"script.src.match(/t=([^&#]*)/)",     // Challenge ID token
	}
	challenged := false
	for _, marker := range challengeMarkers {
		if strings.Contains(strings.ToLower(snippet), strings.ToLower(marker)) {
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
		fmt.Println("🚧 Blocked due to JavaScript/Challenge in body")
	} else if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		fmt.Printf("🔁 Detected redirect to: %s\n", resp.Header.Get("Location"))
	}
	return false
}

