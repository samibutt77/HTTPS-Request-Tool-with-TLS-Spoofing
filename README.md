# HTTPS-Request-Tool-with-TLS-Spoofing
A high-performance, Go-based HTTP client designed to bypass anti-bot mechanisms like Cloudflare, PerimeterX, and Akamai. It mimics real browser behavior through TLS fingerprint spoofing, header mimicry, and intelligent proxy usage.

# Features
🔐 TLS Fingerprint Spoofing using uTLS (Chrome, Firefox, Safari presets)

🌐 Proxy Rotation (HTTP and SOCKS5 with optional authentication)

🎭 Browser Header Spoofing

🍪 Cookie Management via in-memory Jar

🧠 Challenge Detection:

 - Status Code 403

 - CAPTCHA markers in HTML

 - Redirect-based challenges

⏱️ Configurable Request Timing (Min/Max Delay)

🔄 Retry logic with alternate TLS profiles per proxy

🔑 JA3 Hash Calculation & Logging (for each fingerprint)

📁 Single proxy file input with fallback (tries both HTTP and SOCKS5 for each)

# Prerequisites
Before running the tool, ensure you have:

Go 1.20+ installed
Install Go from https://golang.org/dl/

Required Go modules (installed automatically using go mod tidy)

Working proxy list in the format:
- user:pass@ip:port

Save it as proxies.txt in the root directory (right now a proxies.txt file is already present, so that can be used.

# Installation & Setup

- git clone https://github.com/your-repo/abtls.git

- cd abtls

- go mod tidy

Ensure your proxy file proxies.txt is placed in the root directory.

# Usage

Run from terminal:

 - go run cmd/abtls/main.go --url "<TARGET_URL>" --profile <chrome|firefox|safari|random> --min-delay 500 --max-delay 3000

Example:

- go run cmd/abtls/main.go --url "https://www.viagogo.com/Concert-Tickets/Alternative-Music/Coldplay-Tickets/E-155741198" --profile safari --min-delay 500 --max-delay 3000


# Options:

Flag	Description

--url	Target URL to request

--profile	TLS profile (chrome, firefox, safari, or random)

--min-delay	Minimum delay between proxies (ms)

--max-delay	Maximum delay between proxies (ms)

# Proxy Behavior

- Tool uses a single list (proxies.txt)

- Each proxy is:

  - Tried first as HTTP

  - If that fails, retried as SOCKS5

- Proxies can include credentials (user:pass@ip:port) or be simple ip:port 

- Randomly shuffled before processing

# Challenge Detection (Body Heuristics)

The tool analyzes response bodies to detect common anti-bot challenges and blocks.

It looks for HTML/text markers that typically indicate bot detection:

 - cf-challenge

 - g-recaptcha

 - "verify you are human"

Responses are categorized as:

 - ✅ Success (200 OK, clean body)

 - 🚫 Blocked (403 status)

 - 🚧 Challenged (CAPTCHA or challenge detected)

# Cookie Handling & Redirect Support
A built-in cookie jar stores cookies in-memory — either per request or per proxy session.

All requests behave like browsers:

 - Send & receive cookies

 - Maintain session state

Optionally follow HTTP redirects (301, 302, etc.), even when leading to challenge or login pages.

The system can be configured to toggle cookie persistence on or off.



# JA3 Fingerprinting

Each TLS fingerprint (ClientHello) is hashed into a JA3 string

If the JA3 isn't already known (in known_JA3.txt), it's appended

Helps build a list of successful fingerprints over time

# Output Sample

<img width="1915" height="848" alt="image (2)" src="https://github.com/user-attachments/assets/8ebc21ed-a323-443d-8694-c2075b496d58" />


  



