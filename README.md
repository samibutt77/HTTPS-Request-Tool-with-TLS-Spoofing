# HTTPS-Request-Tool-with-TLS-Spoofing
A high-performance, Go-based HTTP client designed to bypass anti-bot mechanisms like Cloudflare, PerimeterX, and Akamai. It mimics real browser behavior through TLS fingerprint spoofing, header mimicry, and intelligent proxy usage.

# Features
🔐 TLS Fingerprint Spoofing using uTLS (Chrome, Firefox, Safari presets)

🌐 Proxy Rotation

🎭 Browser Header Spoofing

🍪 Cookie Management via in-memory Jar

🧠 Challenge Detection:

 - Status Code 403

 - CAPTCHA markers in HTML

 - Redirect-based challenges

⏱️ Configurable Request Timing (Min/Max Delay)

🔄 Retry logic with alternate TLS profiles per proxy

🔑 JA3 Hash Calculation & Logging (for each fingerprint)


# Prerequisites
Before running the tool, ensure you have:

Go 1.20+ installed
Install Go from https://golang.org/dl/

Git Installed
Install Git from https://git-scm.com/downloads/win

Required Go modules (installed automatically using go mod tidy, see in Installation & Setup)

Working proxy list in the format:
- user:pass@ip:port

Save it as proxies.txt in the root directory (right now a proxies.txt file is already present, so that can be used.

# Installation & Setup

- git clone https://github.com/samibutt77/HTTPS-Request-Tool-with-TLS-Spoofing.git

- cd HTTPS-Request-Tool-with-TLS-Spoofing

- cd abtls

- go mod tidy

Ensure your proxy file proxies.txt is placed in the root directory.

# Usage

Run from terminal:

 - go run cmd/abtls/main.go --url "<TARGET_URL>" --profile <chrome|firefox|safari|random> --min-delay 500 --max-delay 3000 --shuffle=<true|false>

 - "random" Profile with shuffle enabled is preferred as I was getting great results on it.

Example:

- go run cmd/abtls/main.go --url "https://www.viagogo.com/Concert-Tickets/Alternative-Music/Coldplay-Tickets/E-155741198" --profile random --min-delay 3000 --max-delay 5000 --shuffle=true

- Enter total number of requests as input.

# Output Sample

<img width="910" height="715" alt="image" src="https://github.com/user-attachments/assets/d5d9f604-bcb8-4dcc-b21d-6f11f841daf5" />


# Additional

- To print known benign JA3, run the command "go run cmd/abtls/main.go --list-ja3"


# Options:

Flag	Description

--url	Target URL to request

--profile	TLS profile (chrome, firefox, safari, or random)

--min-delay	Minimum delay between proxies (ms)

--max-delay	Maximum delay between proxies (ms)

--shuffle Determines if proxies are shuffled or the order in proxies.txt file is followed

# Proxy Behavior

Proxies are loaded from proxies.txt in the format:
 - user:pass@ip:port

Optional shuffle proxies with --shuffle=true (default is in order according to the order in proxies.txt file)

Each proxy is tried using HTTP:

 - If it returns 200 OK, the tool keeps using it until it fails.

 - If it fails or gets blocked, it moves to the next proxy.

Stops after reaching the user-defined total request limit.

Successful proxy-profile combinations (on 200 OK) are stored once in successful_combos.txt.

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






  



