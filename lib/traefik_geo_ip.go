package lib

import (
	"net/http"
)

// TraefikGeoIP is a middleware that put ip in header.
type TraefikGeoIP struct {
	Next                      http.Handler
	Name                      string
	PreferXForwardedForHeader bool
}

func (mw *TraefikGeoIP) ServeHTTP(reqWr http.ResponseWriter, req *http.Request) {
	ipStr := getClientIP(req, mw.PreferXForwardedForHeader)
	req.Header.Set(IPAddressHeader, ipStr)
	mw.Next.ServeHTTP(reqWr, req)
}
