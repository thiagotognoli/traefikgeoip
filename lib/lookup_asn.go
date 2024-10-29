package lib

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/thiagotognoli/traefikgeoip/geoip2"
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

// NewLookupAsn Create a new Lookup.
func NewLookupAsn(dbPath, name string) (LookupGeoIPAsn, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, fmt.Errorf("asn DB not found: db=%s, name=%s, err=%w", dbPath, name, err)
	}
	var lookupAsn LookupGeoIPAsn

	rdr, err := geoip2.NewASNReaderFromFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("asn lookup DB is not initialized: db=%s, name=%s, err=%w", dbPath, name, err)
	}
	lookupAsn = CreateAsnDBLookup(rdr)
	// log.Printf("[geoip2] ASN lookup DB initialized: db=%s, name=%s, lookup=%v", dbPath, name, lookupAsn)
	return lookupAsn, nil
}
