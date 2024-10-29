// Package traefikgeoip is a Traefik plugin for Maxmind GeoIP2.
package traefikgeoip

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	lookupAsn     LookupGeoIPAsn
	lookupCity    LookupGeoIPCity
	lookupCountry LookupGeoIPCountry
)

// ResetLookup reset lookup function.
func ResetLookup() {
	lookupAsn = nil
	lookupCity = nil
	lookupCountry = nil
}

// Config the plugin configuration.
type Config struct {
	CityDBPath                string `json:"cityDbPath,omitempty"`
	AsnDBPath                 string `json:"asnDbPath,omitempty"`
	CountryDBPath             string `json:"countryDbPath,omitempty"`
	PreferXForwardedForHeader bool
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		// CityDBPath: DefaultDBPath,
	}
}

// TraefikGeoIP2 a traefik geoip2 plugin.
type TraefikGeoIP2 struct {
	next                      http.Handler
	name                      string
	preferXForwardedForHeader bool
}

// New created a new TraefikGeoIP plugin.
func New(_ context.Context, next http.Handler, cfg *Config, name string) (http.Handler, error) {
	if cfg.CityDBPath != "" {
		if _, err := os.Stat(cfg.CityDBPath); err != nil {
			log.Printf("[geoip2] City DB not found: db=%s, name=%s, err=%v", cfg.CityDBPath, name, err)
			return &TraefikGeoIP2{
				next:                      next,
				name:                      name,
				preferXForwardedForHeader: cfg.PreferXForwardedForHeader,
			}, nil
		}
		lookupCity, _ = NewLookupCity(lookupCity, cfg.CityDBPath, name)
	} else if cfg.CountryDBPath != "" {
		if _, err := os.Stat(cfg.CountryDBPath); err != nil {
			log.Printf("[geoip2] Country DB not found: db=%s, name=%s, err=%v", cfg.CountryDBPath, name, err)
			return &TraefikGeoIP2{
				next:                      next,
				name:                      name,
				preferXForwardedForHeader: cfg.PreferXForwardedForHeader,
			}, nil
		}
		lookupCountry, _ = NewLookupCountry(lookupCountry, cfg.CountryDBPath, name)
	}
	if cfg.AsnDBPath != "" {
		if _, err := os.Stat(cfg.AsnDBPath); err != nil {
			log.Printf("[geoip2] ASN DB not found: db=%s, name=%s, err=%v", cfg.AsnDBPath, name, err)
			return &TraefikGeoIP2{
				next:                      next,
				name:                      name,
				preferXForwardedForHeader: cfg.PreferXForwardedForHeader,
			}, nil
		}
		lookupAsn, _ = NewLookupAsn(lookupAsn, cfg.AsnDBPath, name)
	}

	return &TraefikGeoIP2{
		next:                      next,
		name:                      name,
		preferXForwardedForHeader: cfg.PreferXForwardedForHeader,
	}, nil
}

func (mw *TraefikGeoIP2) ServeHTTP(reqWr http.ResponseWriter, req *http.Request) {
	ipStr := getClientIP(req, mw.preferXForwardedForHeader)
	req.Header.Set(IPAddressHeader, ipStr)
	switch {
	case lookupCity != nil:
		res, err := lookupCity(net.ParseIP(ipStr))
		if err != nil {
			log.Printf("[geoip2] Unable to find City: ip=%s, err=%v", ipStr, err)
			req.Header.Set(CountryHeader, Unknown)
			req.Header.Set(CountryCodeHeader, Unknown)
			req.Header.Set(RegionHeader, Unknown)
			req.Header.Set(RegionCodeHeader, Unknown)
			req.Header.Set(CityHeader, Unknown)
			req.Header.Set(LatitudeHeader, Unknown)
			req.Header.Set(LongitudeHeader, Unknown)
			req.Header.Set(AccuracyRadiusHeader, Unknown)
			req.Header.Set(GeohashHeader, Unknown)
			req.Header.Set(PostalCodeHeader, Unknown)
		} else {
			req.Header.Set(CountryHeader, res.country)
			req.Header.Set(CountryCodeHeader, res.countryCode)
			req.Header.Set(RegionHeader, res.region)
			req.Header.Set(RegionCodeHeader, res.regionCode)
			req.Header.Set(CityHeader, res.city)
			req.Header.Set(LatitudeHeader, res.latitude)
			req.Header.Set(LongitudeHeader, res.longitude)
			req.Header.Set(AccuracyRadiusHeader, res.accuracyRadius)
			req.Header.Set(GeohashHeader, res.geohash)
			req.Header.Set(PostalCodeHeader, res.postalCode)
		}
		if lookupAsn != nil {
			res, err := lookupAsn(net.ParseIP(ipStr))
			if err != nil {
				log.Printf("[geoip2] Unable to find ASN: ip=%s, err=%v", ipStr, err)
				req.Header.Set(ASNSystemNumberHeader, res.number)
				req.Header.Set(ASNOrganizationHeader, res.organization)
			} else {
				req.Header.Set(ASNSystemNumberHeader, Unknown)
				req.Header.Set(ASNOrganizationHeader, Unknown)
			}
		}
	case lookupCountry != nil:
		res, err := lookupCountry(net.ParseIP(ipStr))
		if err != nil {
			log.Printf("[geoip2] Unable to find Country: ip=%s, err=%v", ipStr, err)
			req.Header.Set(CountryHeader, Unknown)
			req.Header.Set(CountryCodeHeader, Unknown)
		} else {
			req.Header.Set(CountryHeader, res.country)
			req.Header.Set(CountryCodeHeader, res.countryCode)
		}
	case lookupAsn != nil:
		res, err := lookupAsn(net.ParseIP(ipStr))
		if err != nil {
			log.Printf("[geoip2] Unable to find ASN: ip=%s, err=%v", ipStr, err)
			req.Header.Set(ASNSystemNumberHeader, res.number)
			req.Header.Set(ASNOrganizationHeader, res.organization)
		} else {
			req.Header.Set(ASNSystemNumberHeader, Unknown)
			req.Header.Set(ASNOrganizationHeader, Unknown)
		}
	default:
		req.Header.Set(CountryHeader, Unknown)
		req.Header.Set(CountryCodeHeader, Unknown)
		req.Header.Set(RegionHeader, Unknown)
		req.Header.Set(RegionCodeHeader, Unknown)
		req.Header.Set(CityHeader, Unknown)
		req.Header.Set(LatitudeHeader, Unknown)
		req.Header.Set(LongitudeHeader, Unknown)
		req.Header.Set(AccuracyRadiusHeader, Unknown)
		req.Header.Set(GeohashHeader, Unknown)
		req.Header.Set(PostalCodeHeader, Unknown)
		req.Header.Set(ASNSystemNumberHeader, Unknown)
		req.Header.Set(ASNOrganizationHeader, Unknown)
	}

	mw.next.ServeHTTP(reqWr, req)
}

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
