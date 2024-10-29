// Package traefikgeoip is a Traefik plugin for Maxmind GeoIP2.
package traefikgeoip

import (
	"context"
	"log"
	"net/http"

	lib "github.com/thiagotognoli/traefikgeoip/lib"
)

// ResetLookup reset lookup function.
func ResetLookup() {
	// lookupAsn = nil
	// lookupCity = nil
	// lookupCountry = nil
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *lib.Config {
	return &lib.Config{
		// CityDBPath: DefaultDBPath,
	}
}

// New created a new TraefikGeoIP plugin.
func New(_ context.Context, next http.Handler, cfg *lib.Config, name string) (http.Handler, error) {
	lookupCity, lookupCountry, lookupAsn, err := factoryLookups(cfg, name)
	if err != nil {
		log.Printf("%s", err.Error())
		return &lib.TraefikGeoIP{
			Next:                      next,
			Name:                      name,
			PreferXForwardedForHeader: cfg.PreferXForwardedForHeader,
		}, nil // err
	}

	switch {
	case lookupCity != nil && lookupAsn != nil:
		return &lib.TraefikGeoIPCityAsn{
			Next:                      next,
			Name:                      name,
			PreferXForwardedForHeader: cfg.PreferXForwardedForHeader,
			LookupAsn:                 lookupAsn,
			LookupCity:                lookupCity,
		}, nil
	case lookupCity != nil:
		return &lib.TraefikGeoIPCity{
			Next:                      next,
			Name:                      name,
			PreferXForwardedForHeader: cfg.PreferXForwardedForHeader,
			LookupCity:                lookupCity,
		}, nil
	case lookupCountry != nil && lookupAsn != nil:
		return &lib.TraefikGeoIPCountryAsn{
			Next:                      next,
			Name:                      name,
			PreferXForwardedForHeader: cfg.PreferXForwardedForHeader,
			LookupAsn:                 lookupAsn,
			LookupCountry:             lookupCountry,
		}, nil
	case lookupCountry != nil:
		return &lib.TraefikGeoIPCountry{
			Next:                      next,
			Name:                      name,
			PreferXForwardedForHeader: cfg.PreferXForwardedForHeader,
			LookupCountry:             lookupCountry,
		}, nil
	case lookupAsn != nil:
		return &lib.TraefikGeoIPAsn{
			Next:                      next,
			Name:                      name,
			PreferXForwardedForHeader: cfg.PreferXForwardedForHeader,
			LookupAsn:                 lookupAsn,
		}, nil
	default:
		return &lib.TraefikGeoIPNotFound{
			Next:                      next,
			Name:                      name,
			PreferXForwardedForHeader: cfg.PreferXForwardedForHeader,
		}, nil // fmt.Errorf("none GeoIP DB configured")
	}
}

func factoryLookups(cfg *lib.Config, name string) (lib.LookupGeoIPCity, lib.LookupGeoIPCountry, lib.LookupGeoIPAsn, error) {
	var lookupCity lib.LookupGeoIPCity
	var lookupCountry lib.LookupGeoIPCountry
	var lookupAsn lib.LookupGeoIPAsn

	if cfg.CityDBPath != "" {
		var err error
		lookupCity, err = lib.NewLookupCity(cfg.CityDBPath, name)
		if err != nil {
			return nil, nil, nil, err
		}
	} else if cfg.CountryDBPath != "" {
		var err error
		lookupCountry, err = lib.NewLookupCountry(cfg.CountryDBPath, name)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	if cfg.AsnDBPath != "" {
		var err error
		lookupAsn, err = lib.NewLookupAsn(cfg.AsnDBPath, name)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return lookupCity, lookupCountry, lookupAsn, nil
}
