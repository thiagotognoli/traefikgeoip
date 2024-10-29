package lib

import (
	"log"
	"net"
	"net/http"
)

// TraefikGeoIPAsn is a middleware that looks up the city of the client IP address from the GeoIP2 database.
type TraefikGeoIPAsn struct {
	Next      http.Handler
	Name      string
	Options   Options
	LookupAsn LookupGeoIPAsn
}

func (mw *TraefikGeoIPAsn) ServeHTTP(reqWr http.ResponseWriter, req *http.Request) {
	ipStr := getClientIP(req, mw.Options)
	req.Header.Set(IPAddressHeader, ipStr)
	res, err := mw.LookupAsn(net.ParseIP(ipStr))
	if err != nil {
		if mw.Options.Debug {
			log.Printf("[geoip2] Unable to find ASN: ip=%s, err=%v", ipStr, err)
		}
		req.Header.Set(ASNSystemNumberHeader, Unknown)
		req.Header.Set(ASNOrganizationHeader, Unknown)
	} else {
		req.Header.Set(ASNSystemNumberHeader, res.number)
		req.Header.Set(ASNOrganizationHeader, res.organization)
	}
	mw.Next.ServeHTTP(reqWr, req)
}
