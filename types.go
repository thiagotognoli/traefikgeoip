package traefikgeoip

// Unknown constant for undefined data.

// DefaultDBPath default GeoIP2 database path.
const DefaultDBPath = "GeoLite2-City.mmdb"

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
	// AccuracyRadiusHeader coord accuracy radius header name.
	AccuracyRadiusHeader = "GeoIP-Accuracy-Radius"
	// GeohashHeader geohash header name.
	GeohashHeader = "GeoIP-Geohash"

	// ASNSystemNumberHeader asn system number header name.
	ASNSystemNumberHeader = "GeoIP-ASN-System-Number"
	// ASNOrganizationHeader asn system organization header name.
	ASNOrganizationHeader = "GeoIP-ASN-Organization"

	// IPAddressHeader up used in geoip header name.
	IPAddressHeader = "GeoIP-IPAddress"
)
