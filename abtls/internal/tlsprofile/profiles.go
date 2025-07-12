package tlsprofile

import (
	"math/rand"

	utls "github.com/refraction-networking/utls"
)

// GetClientHello returns the ClientHelloSpec and the profile label used.
func GetClientHello(profile string) (utls.ClientHelloSpec, string) {
	switch profile {
	case "chrome":
		spec, _ := utls.UTLSIdToSpec(utls.HelloChrome_Auto)
		return spec, "chrome"
	case "firefox":
		spec, _ := utls.UTLSIdToSpec(utls.HelloFirefox_Auto)
		return spec, "firefox"
	case "safari":
		spec, _ := utls.UTLSIdToSpec(utls.HelloSafari_Auto)
		return spec, "safari"
	case "random":
		// Randomly pick one of the presets
		presets := []struct {
			ID     utls.ClientHelloID
			Label  string
		}{
			{utls.HelloChrome_Auto, "chrome"},
			{utls.HelloFirefox_Auto, "firefox"},
			{utls.HelloSafari_Auto, "safari"},
		}
		pick := presets[rand.Intn(len(presets))]
		spec, _ := utls.UTLSIdToSpec(pick.ID)
		return spec, pick.Label
	default:
		spec, _ := utls.UTLSIdToSpec(utls.HelloRandomized)
		return spec, "randomized"
	}
}
