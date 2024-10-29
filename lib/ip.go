package lib

import (
	"net"
	"net/http"
	"strings"
)

func getClientIP(req *http.Request, preferXForwardedForHeader bool) string {
	if preferXForwardedForHeader {
		// Check X-Forwarded-For header first
		forwardedFor := req.Header.Get("X-Forwarded-For")
		if forwardedFor != "" {
			ips := strings.Split(forwardedFor, ",")
			return strings.TrimSpace(ips[0])
		}
	}

	// If X-Forwarded-For is not present or retrieval is not enabled, fallback to RemoteAddr
	remoteAddr := req.RemoteAddr
	tmp, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		remoteAddr = tmp
	}
	return remoteAddr
}
