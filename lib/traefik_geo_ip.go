package lib

import (
	"net/http"
)

// TraefikGeoIP is a middleware that put ip in header.
type TraefikGeoIP struct {
	Next    http.Handler
	Name    string
	Options Options
}

func (mw *TraefikGeoIP) ServeHTTP(reqWr http.ResponseWriter, req *http.Request) {
	ipStr := getClientIP(req, mw.Options)
	req.Header.Set(IPAddressHeader, ipStr)
	mw.Next.ServeHTTP(reqWr, req)
}
