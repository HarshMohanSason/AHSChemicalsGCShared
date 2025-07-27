package utils

import (
	"net"
	"net/http"
	"strings"
)

// GetIp attempts to retrieve the IP address of the client making the HTTP request.
//
// This function checks for commonly used headers set by proxies and load balancers
// to extract the original client IP address in the following order:
//   1. "X-Forwarded-For": May contain a comma-separated list of IPs. The first one is typically the client IP.
//   2. "X-Real-Ip": Set by some reverse proxies to indicate the client’s IP.
//   3. request.RemoteAddr: Falls back to the address directly from the TCP connection.
//
// It returns the first valid, parsed IP address found in these sources, or an empty string if none is found.
//
// Parameters:
//   - request (*http.Request): The HTTP request from which to extract the client IP.
//
// Returns:
//   - string: The client's IP address, or an empty string if it cannot be determined.
//
func GetIp(request *http.Request) string {
	xff := request.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		ip := strings.TrimSpace(ips[0])
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}

	xrip := request.Header.Get("X-Real-Ip")
	if xrip != "" {
		if parsed := net.ParseIP(xrip); parsed != nil {
			return parsed.String()
		}
	}

	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil {
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}

	return ""
}