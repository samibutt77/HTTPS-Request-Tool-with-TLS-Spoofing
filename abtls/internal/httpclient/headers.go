package httpclient

//import "net/http"

type Header struct {
	Key   string
	Value string
}

func GetOrderedHeadersForProfile(profile string) []Header {
	switch profile {
	case "chrome":
		return []Header{
			{"Sec-Ch-Ua", `"Chromium";v="114", "Not.A/Brand";v="8", "Google Chrome";v="114"`},
			{"Sec-Ch-Ua-Mobile", "?0"},
			{"Sec-Ch-Ua-Platform", `"Windows"`},
			{"Upgrade-Insecure-Requests", "1"},
			{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
			{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
			{"Sec-Fetch-Site", "none"},
			{"Sec-Fetch-Mode", "navigate"},
			{"Sec-Fetch-User", "?1"},
			{"Sec-Fetch-Dest", "document"},
			{"Accept-Encoding", "gzip, deflate, br"},
			{"Accept-Language", "en-US,en;q=0.9"},
			{"Referer", "https://www.google.com/"},
			{"Origin", "https://www.google.com"},
		}

	case "firefox":
		return []Header{
			{"Upgrade-Insecure-Requests", "1"},
			{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0"},
			{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
			{"Accept-Language", "en-US,en;q=0.5"},
			{"Accept-Encoding", "gzip, deflate, br"},
			{"Sec-Fetch-Site", "none"},
			{"Sec-Fetch-Mode", "navigate"},
			{"Sec-Fetch-User", "?1"},
			{"Sec-Fetch-Dest", "document"},
			{"Referer", "https://www.google.com/"},
			{"Origin", "https://www.google.com"},
		}

	case "safari":
		return []Header{
			{"Upgrade-Insecure-Requests", "1"},
			{"User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15"},
			{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
			{"Accept-Language", "en-us"},
			{"Accept-Encoding", "gzip, deflate, br"},
			{"Referer", "https://www.google.com/"},
			{"Origin", "https://www.google.com"},
		}

	default:
		return nil
	}
}
