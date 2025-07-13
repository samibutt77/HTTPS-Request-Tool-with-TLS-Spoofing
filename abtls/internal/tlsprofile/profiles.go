package tlsprofile

import (
	"math/rand"

	utls "github.com/refraction-networking/utls"
)

// GetClientHello returns the ClientHelloSpec and the profile label used.
func GetClientHello(profile string) (utls.ClientHelloSpec, string) {
	addHTTP2ALPN := func(spec utls.ClientHelloSpec) utls.ClientHelloSpec {
		// Append or replace ALPNExtension with h2 and http/1.1
		spec.Extensions = append(spec.Extensions, &utls.ALPNExtension{
			AlpnProtocols: []string{"h2", "http/1.1"},
		})
		return spec
	}

	switch profile {
	case "chrome":
		spec, _ := utls.UTLSIdToSpec(utls.HelloChrome_Auto)
		spec = addHTTP2ALPN(spec)
		return spec, "chrome"

	case "firefox":
		spec, _ := utls.UTLSIdToSpec(utls.HelloFirefox_Auto)
		spec = addHTTP2ALPN(spec)
		return spec, "firefox"

	case "safari":
		// Safari typically doesn't negotiate HTTP/2 explicitly in some versions
		spec, _ := utls.UTLSIdToSpec(utls.HelloSafari_Auto)
		return spec, "safari"

	case "random":
		presets := []struct {
			ID    utls.ClientHelloID
			Label string
		}{
			{utls.HelloChrome_Auto, "chrome"},
			{utls.HelloFirefox_Auto, "firefox"},
			{utls.HelloSafari_Auto, "safari"},
		}
		pick := presets[rand.Intn(len(presets))]
		spec, _ := utls.UTLSIdToSpec(pick.ID)
		if pick.Label == "chrome" || pick.Label == "firefox" {
			spec = addHTTP2ALPN(spec)
		}
		return spec, pick.Label

	default:
		spec, _ := utls.UTLSIdToSpec(utls.HelloRandomized)
		return spec, "randomized"
	}
}
