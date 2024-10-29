package lib

import (
	"log"
	"net"
	"net/http"
)

// TraefikGeoIPCountryAsn is a middleware that looks up the city of the client IP address from the GeoIP2 database.
type TraefikGeoIPCountryAsn struct {
	Next                      http.Handler
	Name                      string
	PreferXForwardedForHeader bool
	LookupAsn                 LookupGeoIPAsn
	LookupCountry             LookupGeoIPCountry
}

func (mw *TraefikGeoIPCountryAsn) ServeHTTP(reqWr http.ResponseWriter, req *http.Request) {
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
	resAsn, err := mw.LookupAsn(net.ParseIP(ipStr))
	if err != nil {
		log.Printf("[geoip2] Unable to find ASN: ip=%s, err=%v", ipStr, err)
		req.Header.Set(ASNSystemNumberHeader, Unknown)
		req.Header.Set(ASNOrganizationHeader, Unknown)
	} else {
		req.Header.Set(ASNSystemNumberHeader, resAsn.number)
		req.Header.Set(ASNOrganizationHeader, resAsn.organization)
	}

	mw.Next.ServeHTTP(reqWr, req)
}
