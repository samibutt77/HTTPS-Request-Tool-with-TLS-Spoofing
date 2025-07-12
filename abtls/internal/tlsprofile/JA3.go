package tlsprofile

import (
	"crypto/md5"
	"encoding/hex" // ✅ Required for hex.EncodeToString
	"fmt"
	"strings"

	utls "github.com/refraction-networking/utls"
)

// ComputeJA3 builds the JA3 string from the ClientHelloSpec and returns its MD5 hash.
func computeJA3Hash(spec utls.ClientHelloSpec) string {
	var b strings.Builder

	// TLS Version (always TLS 1.2 for uTLS)
	b.WriteString("771")

	// Cipher Suites
	b.WriteString(",")
	for i, cs := range spec.CipherSuites {
		b.WriteString(fmt.Sprintf("%d", cs))
		if i < len(spec.CipherSuites)-1 {
			b.WriteString("-")
		}
	}

	// Extensions
	b.WriteString(",")
	var extIDs []int
	for _, ext := range spec.Extensions {
		switch ext.(type) {
		case *utls.SNIExtension:
			extIDs = append(extIDs, 0)
		case *utls.SupportedCurvesExtension:
			extIDs = append(extIDs, 10)
		case *utls.SupportedPointsExtension:
			extIDs = append(extIDs, 11)
		case *utls.SignatureAlgorithmsExtension:
			extIDs = append(extIDs, 13)
		case *utls.ALPNExtension:
			extIDs = append(extIDs, 16)
		case *utls.StatusRequestExtension:
			extIDs = append(extIDs, 5)
		case *utls.SCTExtension:
			extIDs = append(extIDs, 18)
		case *utls.RenegotiationInfoExtension:
			extIDs = append(extIDs, 65281)
		case *utls.KeyShareExtension:
			extIDs = append(extIDs, 51)
		case *utls.PSKKeyExchangeModesExtension:
			extIDs = append(extIDs, 45)
		case *utls.SupportedVersionsExtension:
			extIDs = append(extIDs, 43)
		case *utls.CookieExtension:
			extIDs = append(extIDs, 44)
		case *utls.SessionTicketExtension:
			extIDs = append(extIDs, 35)
		}
	}
	for i, id := range extIDs {
		b.WriteString(fmt.Sprintf("%d", id))
		if i < len(extIDs)-1 {
			b.WriteString("-")
		}
	}

	// Elliptic Curves
	b.WriteString(",")
	var curves []int
	for _, ext := range spec.Extensions {
		if e, ok := ext.(*utls.SupportedCurvesExtension); ok {
			for _, c := range e.Curves {
				curves = append(curves, int(c))
			}
		}
	}
	for i, c := range curves {
		b.WriteString(fmt.Sprintf("%d", c))
		if i < len(curves)-1 {
			b.WriteString("-")
		}
	}

	// EC Point Formats
	b.WriteString(",")
	var points []int
	for _, ext := range spec.Extensions {
		if e, ok := ext.(*utls.SupportedPointsExtension); ok {
			for _, p := range e.SupportedPoints {
				points = append(points, int(p))
			}
		}
	}
	for i, p := range points {
		b.WriteString(fmt.Sprintf("%d", p))
		if i < len(points)-1 {
			b.WriteString("-")
		}
	}

	// Final MD5 hash
	sum := md5.Sum([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}
