package httpclient

import "net/http"

func GetHeadersForProfile(profile string) http.Header {
	headers := http.Header{}

	switch profile {
	case "chrome":
		headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
		headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		headers.Set("Accept-Language", "en-US,en;q=0.9")
		headers.Set("Accept-Encoding", "gzip, deflate, br")
		headers.Set("Sec-Fetch-Dest", "document")
		headers.Set("Sec-Fetch-Mode", "navigate")
		headers.Set("Sec-Fetch-Site", "none")
		headers.Set("Sec-Fetch-User", "?1")
		headers.Set("Upgrade-Insecure-Requests", "1")

	case "firefox":
		headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0")
		headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		headers.Set("Accept-Language", "en-US,en;q=0.5")
		headers.Set("Accept-Encoding", "gzip, deflate, br")
		headers.Set("Upgrade-Insecure-Requests", "1")
		headers.Set("Sec-Fetch-Dest", "document")
		headers.Set("Sec-Fetch-Mode", "navigate")
		headers.Set("Sec-Fetch-Site", "none")
		headers.Set("Sec-Fetch-User", "?1")

	case "safari":
		headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15")
		headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		headers.Set("Accept-Language", "en-us")
		headers.Set("Accept-Encoding", "gzip, deflate, br")
		headers.Set("Upgrade-Insecure-Requests", "1")
	}

	return headers
}
