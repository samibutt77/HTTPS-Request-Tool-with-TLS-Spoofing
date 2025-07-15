package main

import (
	"fmt"
	"os"

	"abtls/internal/config"
	"abtls/internal/runner"
)

func main() {
	cfg := config.ParseFlags() // ✅ Parse all flags from config.go

	// Handle --list-ja3 directly
	if cfg.ListJA3 {
		data, err := os.ReadFile("known_JA3.txt")
		if err != nil {
			fmt.Println("❌ Error reading known_JA3.txt:", err)
			return
		}
		fmt.Println("📋 Known JA3 Hashes:")
		fmt.Println(string(data))
		return
	}

	// Validate required URL
	if cfg.URL == "" {
		fmt.Println("❌ Error: --url is required")
		return
	}

	// Start runner
	runner.Run(cfg)
}
