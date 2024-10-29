package lib

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/thiagotognoli/traefikgeoip/geoip2"
)

// GeoIPCityResult in memory, this should have between 126 and 180 bytes. On average, consider 150 bytes.
type GeoIPCityResult struct {
	country        string
	countryCode    string
	region         string
	regionCode     string
	city           string
	latitude       string
	longitude      string
	accuracyRadius string
	geohash        string
	postalCode     string
}

const kmToMeters = 1000

// LookupGeoIPCity LookupGeoIP.
type LookupGeoIPCity func(ip net.IP) (*GeoIPCityResult, error)

// CreateCityDBLookup CreateCityDBLookup.
func CreateCityDBLookup(rdr *geoip2.CityReader) LookupGeoIPCity {
	return func(ip net.IP) (*GeoIPCityResult, error) {
		rec, err := rdr.Lookup(ip)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		returnVal := GeoIPCityResult{
			country:        Unknown,
			countryCode:    rec.Country.ISOCode,
			region:         Unknown,
			regionCode:     Unknown,
			city:           Unknown,
			postalCode:     rec.Postal.Code,
			latitude:       strconv.FormatFloat(rec.Location.Latitude, 'f', -1, 64),
			longitude:      strconv.FormatFloat(rec.Location.Longitude, 'f', -1, 64),
			accuracyRadius: strconv.Itoa(int(rec.Location.AccuracyRadius) * kmToMeters),
			geohash:        EncodeGeoHash(rec.Location.Latitude, rec.Location.Longitude),
		}
		if country, ok := rec.Country.Names["en"]; ok {
			returnVal.country = country
		}
		if city, ok := rec.City.Names["en"]; ok {
			returnVal.city = city
		}
		if rec.Subdivisions != nil {
			if region, ok := rec.Subdivisions[0].Names["en"]; ok {
				returnVal.region = region
			}
			returnVal.regionCode = rec.Subdivisions[0].ISOCode
		}
		return &returnVal, nil
	}
}

// NewLookupCity Create a new Lookup.
func NewLookupCity(dbPath, name string) (LookupGeoIPCity, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, fmt.Errorf("city DB not found: db=%s, name=%s, err=%w", dbPath, name, err)
	}
	var lookupCity LookupGeoIPCity

	rdr, err := geoip2.NewCityReaderFromFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("city lookup DB is not initialized: db=%s, name=%s, err=%w", dbPath, name, err)
	}
	lookupCity = CreateCityDBLookup(rdr)
	// log.Printf("[geoip2] City lookup DB initialized: db=%s, name=%s, lookup=%v", dbPath, name, lookupCity)
	return lookupCity, nil
}
