package lib

import (
	"log"
	"net"
	"net/http"
)

// TraefikGeoIPCountry is a middleware that looks up the city of the client IP address from the GeoIP2 database.
type TraefikGeoIPCountry struct {
	Next                      http.Handler
	Name                      string
	PreferXForwardedForHeader bool
	LookupCountry             LookupGeoIPCountry
}

func (mw *TraefikGeoIPCountry) ServeHTTP(reqWr http.ResponseWriter, req *http.Request) {
	ipStr := getClientIP(req, mw.PreferXForwardedForHeader)
	req.Header.Set(IPAddressHeader, ipStr)
	res, err := mw.LookupCountry(net.ParseIP(ipStr))
	if err != nil {
		log.Printf("[geoip2] Unable to find Country: ip=%s, err=%v", ipStr, err)
		req.Header.Set(CountryHeader, Unknown)
		req.Header.Set(CountryCodeHeader, Unknown)
	} else {
		req.Header.Set(CountryHeader, res.country)
		req.Header.Set(CountryCodeHeader, res.countryCode)
	}
	mw.Next.ServeHTTP(reqWr, req)
}
