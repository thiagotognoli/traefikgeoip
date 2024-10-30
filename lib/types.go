// Package lib package contains traefikgeoip implementations.
package lib

import "net/http"

// TraefikGeoIPBase is a base middleware that looks client IP address from the GeoIP2 database.
type TraefikGeoIPBase struct {
	Next    http.Handler
	Name    string
	Options Options
}

// TraefikGeoIPNotFound is a middleware that do nothing.
type TraefikGeoIPNotFound struct {
	Next    http.Handler
	Name    string
	Options Options
}

func (mw *TraefikGeoIPNotFound) ServeHTTP(reqWr http.ResponseWriter, req *http.Request) {
	mw.Next.ServeHTTP(reqWr, req)
}

// Options the plugin options.
type Options struct {
	PreferXForwardedForHeader bool
	IPHeader                  string `json:"ipHeader,omitempty"`
	FailInError               bool   `json:"failInError,omitempty"`
	Debug                     bool   `json:"debug,omitempty"`
	LightMode                 bool   `json:"lightMode,omitempty"`
	Iso88591                  bool   `json:"iso88591,omitempty"`
}

// Config the plugin configuration.
type Config struct {
	CityDBPath                string `json:"cityDbPath,omitempty"`
	AsnDBPath                 string `json:"asnDbPath,omitempty"`
	CountryDBPath             string `json:"countryDbPath,omitempty"`
	PreferXForwardedForHeader bool
	IPHeader                  string `json:"ipHeader,omitempty"`
	FailInError               bool   `json:"failInError,omitempty"`
	Debug                     bool   `json:"debug,omitempty"`
	LightMode                 bool   `json:"lightMode,omitempty"`
	Iso88591                  bool   `json:"iso88591,omitempty"`
}

// ConfigToOptions converts the plugin configuration to plugin options.
func ConfigToOptions(config *Config) Options {
	return Options{
		PreferXForwardedForHeader: config.PreferXForwardedForHeader,
		IPHeader:                  config.IPHeader,
		FailInError:               config.FailInError,
		Debug:                     config.Debug,
		LightMode:                 config.LightMode,
		Iso88591:                  config.Iso88591,
	}
}

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
