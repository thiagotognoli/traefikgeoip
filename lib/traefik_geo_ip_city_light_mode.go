package lib

import (
	"log"
	"net"
	"net/http"
)

// TraefikGeoIPCityLightMode is a middleware that looks up the city of the client IP address from the GeoIP2 database.
type TraefikGeoIPCityLightMode struct {
	Next       http.Handler
	Name       string
	Options    Options
	LookupCity LookupGeoIPCity
}

func (mw *TraefikGeoIPCityLightMode) ServeHTTP(reqWr http.ResponseWriter, req *http.Request) {
	ipStr := getClientIP(req, mw.Options)
	req.Header.Set(IPAddressHeader, ipStr)
	res, err := mw.LookupCity(net.ParseIP(ipStr))
	if err != nil {
		if mw.Options.Debug {
			log.Printf("[geoip2] Unable to find City: ip=%s, err=%v", ipStr, err)
		}
		req.Header.Set(CountryCodeHeader, Unknown)
		req.Header.Set(RegionCodeHeader, Unknown)
		req.Header.Set(CityHeader, Unknown)
		req.Header.Set(LatitudeHeader, Unknown)
		req.Header.Set(LongitudeHeader, Unknown)
		req.Header.Set(AccuracyRadiusHeader, Unknown)
	} else {
		req.Header.Set(CountryCodeHeader, res.countryCode)
		req.Header.Set(RegionCodeHeader, res.regionCode)
		req.Header.Set(CityHeader, res.city)
		req.Header.Set(LatitudeHeader, res.latitude)
		req.Header.Set(LongitudeHeader, res.longitude)
		req.Header.Set(AccuracyRadiusHeader, res.accuracyRadius)
	}

	mw.Next.ServeHTTP(reqWr, req)
}
