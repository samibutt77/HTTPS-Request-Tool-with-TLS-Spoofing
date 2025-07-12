package config

import (
	"flag"
)

type Config struct {
	URL        string
	ProxyFile  string
	TLSProfile string
	MinDelay   int
	MaxDelay   int
	ListJA3    bool
}

func ParseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.URL, "url", "", "Target URL")
	flag.StringVar(&cfg.ProxyFile, "proxy-file", "proxies.txt", "Path to proxy list")
	flag.StringVar(&cfg.TLSProfile, "profile", "chrome", "TLS profile: chrome/firefox/safari/random")
	flag.IntVar(&cfg.MinDelay, "min-delay", 500, "Minimum delay between requests (in ms)")
	flag.IntVar(&cfg.MaxDelay, "max-delay", 3000, "Maximum delay between requests (in ms)")
	

	return cfg
}
