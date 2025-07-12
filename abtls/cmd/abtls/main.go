package main

import (
	"flag"
	"fmt"
	"os"

	"abtls/internal/config"
	"abtls/internal/runner"
)

func main() {
	// Define config flags directly
	cfg := &config.Config{}
	flag.StringVar(&cfg.URL, "url", "", "Target URL")
	flag.StringVar(&cfg.ProxyFile, "proxy-file", "proxies.txt", "Path to proxy list")
	flag.StringVar(&cfg.TLSProfile, "profile", "chrome", "TLS profile: chrome/firefox/random")
	flag.IntVar(&cfg.MinDelay, "min-delay", 500, "Minimum delay between requests (in ms)")
	flag.IntVar(&cfg.MaxDelay, "max-delay", 3000, "Maximum delay between requests (in ms)")

	// Add JA3 listing flag
	listJA3 := flag.Bool("list-ja3", false, "List known JA3 hashes from known_JA3.txt")

	flag.Parse() // ✅ Parse all at once

	// Show known JA3 hashes and exit
	if *listJA3 {
		data, err := os.ReadFile("known_JA3.txt")
		if err != nil {
			fmt.Println("❌ Error reading known_JA3.txt:", err)
			return
		}
		fmt.Println("📋 Known JA3 Hashes:")
		fmt.Println(string(data))
		return
	}

	// Validate URL
	if cfg.URL == "" {
		fmt.Println("❌ Error: --url is required")
		return
	}

	// Run main logic
	runner.Run(cfg)
}
