package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetIp(request *http.Request) string {
	// Check X-Forwarded-For header
	xff := request.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		ip := strings.TrimSpace(ips[0])
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}

	// Check X-Real-IP header
	xrip := request.Header.Get("X-Real-Ip")
	if xrip != "" {
		if parsed := net.ParseIP(xrip); parsed != nil {
			return parsed.String()
		}
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil {
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}

	// No IP found
	return ""
}
