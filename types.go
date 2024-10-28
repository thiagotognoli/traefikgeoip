package traefikgeoip

import (
	"fmt"
	"net"

	"github.com/thiagotognoli/traefikgeoip/geoip2"
)

// Unknown constant for undefined data.

// DefaultDBPath default GeoIP2 database path.
const DefaultDBPath = "GeoLite2-Country.mmdb"

// const Unknown = "XX"
// const (
// 	// CountryHeader country header name.
// 	CountryHeader = "X-GeoIP2-Country"
// 	// RegionHeader region header name.
// 	RegionHeader = "X-GeoIP2-Region"
// 	// CityHeader city header name.
// 	CityHeader = "X-GeoIP2-City"
// 	// IPAddressHeader city header name.
// 	IPAddressHeader = "X-GeoIP2-IPAddress"
// )

const (
	// Unknown constant for undefined data.
	Unknown = "XX"
	// ContinentHeader country header name.
	ContinentHeader = "GeoIP-Continent"
	// ContinentCodeHeader country code header name.
	ContinentCodeHeader = "GeoIP-Continent-Code"
	// CountryHeader country header name.
	CountryHeader = "GeoIP-Country"
	// CountryCodeHeader country code header name.
	CountryCodeHeader = "GeoIP-Country-Code"
	// RegionHeader region header name.
	RegionHeader = "GeoIP-Region"
	// RegionCodeHeader region code header name.
	RegionCodeHeader = "GeoIP-Region-Code"
	// CityHeader city header name.
	CityHeader = "GeoIP-City"
	// PostalCodeHeader city header name.
	PostalCodeHeader = "GeoIP-Postal-Code"

	// LatitudeHeader latitude header name.
	LatitudeHeader = "GeoIP-Latitude"
	// LongitudeHeader longitude header name.
	LongitudeHeader = "GeoIP-Longitude"
	// GeohashHeader geohash header name.
	GeohashHeader = "GeoIP-Geohash"

	// ASNSystemNumberHeader asn system number header name.
	ASNSystemNumberHeader = "GeoIP-ASN-System-Number"
	// ASNOrganizationHeader asn system organization header name.
	ASNOrganizationHeader = "GeoIP-ASN-Organization"

	// IPAddressHeader up used in geoip header name.
	IPAddressHeader = "GeoIP-IPAddress"
)

// GeoIPResult GeoIPResult.
type GeoIPResult struct {
	country string
	region  string
	city    string
}

// LookupGeoIP2 LookupGeoIP2.
type LookupGeoIP2 func(ip net.IP) (*GeoIPResult, error)

// CreateCityDBLookup CreateCityDBLookup.
func CreateCityDBLookup(rdr *geoip2.CityReader) LookupGeoIP2 {
	return func(ip net.IP) (*GeoIPResult, error) {
		rec, err := rdr.Lookup(ip)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		retval := GeoIPResult{
			country: rec.Country.ISOCode,
			region:  Unknown,
			city:    Unknown,
		}
		if city, ok := rec.City.Names["en"]; ok {
			retval.city = city
		}
		if rec.Subdivisions != nil {
			retval.region = rec.Subdivisions[0].ISOCode
		}
		return &retval, nil
	}
}

// CreateCountryDBLookup CreateCountryDBLookup.
func CreateCountryDBLookup(rdr *geoip2.CountryReader) LookupGeoIP2 {
	return func(ip net.IP) (*GeoIPResult, error) {
		rec, err := rdr.Lookup(ip)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		retval := GeoIPResult{
			country: rec.Country.ISOCode,
			region:  Unknown,
			city:    Unknown,
		}
		return &retval, nil
	}
}
