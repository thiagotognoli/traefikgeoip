package traefikgeoip

import (
	"fmt"
	"log"
	"net"

	"github.com/thiagotognoli/traefikgeoip/geoip2"
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

// NewLookupCountry Create a new Lookup.
func NewLookupCountry(lookupCountry LookupGeoIPCountry, dbPath, name string) (LookupGeoIPCountry, error) {
	// var lookupCountry LookupGeoIPCountry
	// if lookupCountry == nil {
	rdr, err := geoip2.NewCountryReaderFromFile(dbPath)
	if err != nil {
		log.Printf("[geoip2] Country lookup DB is not initialized: db=%s, name=%s, err=%v", dbPath, name, err)
	} else {
		lookupCountry = CreateCountryDBLookup(rdr)
		log.Printf("[geoip2] Country lookup DB initialized: db=%s, name=%s, lookup=%v", dbPath, name, lookupCountry)
	}
	// }

	return lookupCountry, nil
}
