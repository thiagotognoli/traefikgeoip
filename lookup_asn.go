package traefikgeoip

import (
	"fmt"
	"log"
	"net"
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
func NewLookupAsn(lookupAsn LookupGeoIPAsn, dbPath, name string) (LookupGeoIPAsn, error) {
	// if lookupAsn == nil {
	rdr, err := geoip2.NewASNReaderFromFile(dbPath)
	if err != nil {
		log.Printf("[geoip2] ASN lookup DB is not initialized: db=%s, name=%s, err=%v", dbPath, name, err)
	} else {
		lookupAsn = CreateAsnDBLookup(rdr)
		log.Printf("[geoip2] ASN lookup DB initialized: db=%s, name=%s, lookup=%v", dbPath, name, lookupAsn)
	}
	// }

	return lookupAsn, nil
}
