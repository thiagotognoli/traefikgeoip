package lib

import (
	"fmt"
	"net"
	"os"

	geoip2 "github.com/thiagotognoli/traefikgeoip/geoip2"
	geoip2_iso88591 "github.com/thiagotognoli/traefikgeoip/geoip2_iso88591"
)

// GeoIPCountryResult in memory, this should have between 126 and 180 bytes. On average, consider 150 bytes.
type GeoIPCountryResult struct {
	country     string
	countryCode string
}

// LookupGeoIPCountry LookupGeoIPCountry.
type LookupGeoIPCountry func(ip net.IP) (*GeoIPCountryResult, error)

// CreateCountryDBLookup CreateCountryDBLookup.
func CreateCountryDBLookup(rdr *geoip2.CountryReader) LookupGeoIPCountry {
	return func(ip net.IP) (*GeoIPCountryResult, error) {
		rec, err := rdr.Lookup(ip)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		returnVal := GeoIPCountryResult{
			country:     Unknown,
			countryCode: rec.Country.ISOCode,
		}
		if country, ok := rec.Country.Names["en"]; ok {
			returnVal.country = country
		}
		return &returnVal, nil
	}
}

// CreateCountryDBLookupIso88591 CreateCountryDBLookup.
func CreateCountryDBLookupIso88591(rdr *geoip2_iso88591.CountryReader) LookupGeoIPCountry {
	return func(ip net.IP) (*GeoIPCountryResult, error) {
		rec, err := rdr.Lookup(ip)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		returnVal := GeoIPCountryResult{
			country:     Unknown,
			countryCode: rec.Country.ISOCode,
		}
		if country, ok := rec.Country.Names["en"]; ok {
			returnVal.country = country
		}
		return &returnVal, nil
	}
}

// NewLookupCountry Create a new Lookup.
func NewLookupCountry(dbPath, name string, iso88591 bool) (LookupGeoIPCountry, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, fmt.Errorf("country DB not found: db=%s, name=%s, err=%w", dbPath, name, err)
	}
	var lookupCountry LookupGeoIPCountry

	if iso88591 {
		rdr, err := geoip2_iso88591.NewCountryReaderFromFile(dbPath)
		if err != nil {
			return nil, fmt.Errorf("country lookup DB is not initialized: db=%s, name=%s, err=%w", dbPath, name, err)
		}
		lookupCountry = CreateCountryDBLookupIso88591(rdr)
	} else {
		rdr, err := geoip2.NewCountryReaderFromFile(dbPath)
		if err != nil {
			return nil, fmt.Errorf("country lookup DB is not initialized: db=%s, name=%s, err=%w", dbPath, name, err)
		}
		lookupCountry = CreateCountryDBLookup(rdr)
	}
	// log.Printf("[geoip2] Country lookup DB initialized: db=%s, name=%s, lookup=%v", dbPath, name, lookupCountry)
	return lookupCountry, nil
}
