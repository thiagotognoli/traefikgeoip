package lib

import (
	"fmt"
	"net"
	"os"
	"strconv"

	geoip2 "github.com/thiagotognoli/traefikgeoip/geoip2"
	geoip2_iso88591 "github.com/thiagotognoli/traefikgeoip/geoip2_iso88591"
)

// GeoIPAsnResult in memory, this should have between 126 and 180 bytes. On average, consider 150 bytes.
type GeoIPAsnResult struct {
	number       string
	organization string
}

// LookupGeoIPAsn LookupGeoIP.
type LookupGeoIPAsn func(ip net.IP) (*GeoIPAsnResult, error)

// CreateAsnDBLookup CreateCountryDBLookup.
func CreateAsnDBLookup(rdr *geoip2.ASNReader) LookupGeoIPAsn {
	return func(ip net.IP) (*GeoIPAsnResult, error) {
		rec, err := rdr.Lookup(ip)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		returnVal := GeoIPAsnResult{
			number:       strconv.Itoa(int(rec.AutonomousSystemNumber)),
			organization: rec.AutonomousSystemOrganization,
		}
		return &returnVal, nil
	}
}

// CreateAsnDBLookupIso88591 CreateCountryDBLookup.
func CreateAsnDBLookupIso88591(rdr *geoip2_iso88591.ASNReader) LookupGeoIPAsn {
	return func(ip net.IP) (*GeoIPAsnResult, error) {
		rec, err := rdr.Lookup(ip)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		returnVal := GeoIPAsnResult{
			number:       strconv.Itoa(int(rec.AutonomousSystemNumber)),
			organization: rec.AutonomousSystemOrganization,
		}
		return &returnVal, nil
	}
}

// NewLookupAsn Create a new Lookup.
func NewLookupAsn(dbPath, name string, iso88591 bool) (LookupGeoIPAsn, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, fmt.Errorf("asn DB not found: db=%s, name=%s, err=%w", dbPath, name, err)
	}
	var lookupAsn LookupGeoIPAsn

	if iso88591 {
		rdr, err := geoip2_iso88591.NewASNReaderFromFile(dbPath)
		if err != nil {
			return nil, fmt.Errorf("asn lookup DB is not initialized: db=%s, name=%s, err=%w", dbPath, name, err)
		}
		lookupAsn = CreateAsnDBLookupIso88591(rdr)
	} else {
		rdr, err := geoip2.NewASNReaderFromFile(dbPath)
		if err != nil {
			return nil, fmt.Errorf("asn lookup DB is not initialized: db=%s, name=%s, err=%w", dbPath, name, err)
		}
		lookupAsn = CreateAsnDBLookup(rdr)
	}
	// log.Printf("[geoip2] ASN lookup DB initialized: db=%s, name=%s, lookup=%v", dbPath, name, lookupAsn)
	return lookupAsn, nil
}
