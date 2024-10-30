package traefikgeoip_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mw "github.com/thiagotognoli/traefikgeoip"
	lmw "github.com/thiagotognoli/traefikgeoip/lib"
)

const (
	ValidIP          = "188.193.88.199"
	ValidAlternateIP = "188.193.88.200"
	ValidIPNoCity    = "20.1.184.61"
)

func TestGeoIPConfig(t *testing.T) {
	mwCfg := mw.CreateConfig()
	// if mw.DefaultDBPath != mwCfg.CityDBPath {
	// 	t.Fatalf("Incorrect path")
	// }

	mwCfg.CityDBPath = "./non-existing"
	mw.ResetLookup()
	_, err := mw.New(context.TODO(), nil, mwCfg, "")
	if err != nil {
		t.Fatalf("Must not fail on missing DB")
	}

	mwCfg.CityDBPath = "justfile"
	_, err = mw.New(context.TODO(), nil, mwCfg, "")
	if err != nil {
		t.Fatalf("Must not fail on invalid DB format")
	}
}

func TestGeoIPBasic(t *testing.T) {
	mwCfg := mw.CreateConfig()
	mwCfg.CityDBPath = "data/tmp/GeoLite2-City.mmdb"

	called := false
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		called = true
	})

	mw.ResetLookup()
	instance, err := mw.New(context.TODO(), next, mwCfg, "traefik-geoip")
	if err != nil {
		t.Fatalf("Error creating %v", err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)

	instance.ServeHTTP(recorder, req)
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("Invalid return code")
	}
	if called != true {
		t.Fatalf("next handler was not called")
	}
}

func TestMissingGeoIPDB(t *testing.T) {
	mwCfg := mw.CreateConfig()
	mwCfg.CityDBPath = "./missing"

	called := false
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { called = true })

	mw.ResetLookup()
	instance, err := mw.New(context.TODO(), next, mwCfg, "traefik-geoip")
	if err != nil {
		t.Fatalf("Error creating %v", err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = "1.2.3.4"

	instance.ServeHTTP(recorder, req)
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("Invalid return code")
	}
	if called != true {
		t.Fatalf("next handler was not called")
	}
	assertHeader(t, req, lmw.CountryHeader, "")
	assertHeader(t, req, lmw.RegionHeader, "")
	assertHeader(t, req, lmw.CityHeader, "")
	assertHeader(t, req, lmw.IPAddressHeader, "1.2.3.4")
}

func TestGeoIPFromRemoteAddr(t *testing.T) {
	mwCfg := mw.CreateConfig()
	mwCfg.CityDBPath = "data/tmp/GeoLite2-City.mmdb"

	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	mw.ResetLookup()
	instance, _ := mw.New(context.TODO(), next, mwCfg, "traefik-geoip")

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = fmt.Sprintf("%s:9999", ValidIP)
	instance.ServeHTTP(httptest.NewRecorder(), req)
	assertHeader(t, req, lmw.CountryCodeHeader, "DE")
	assertHeader(t, req, lmw.RegionCodeHeader, "BY")
	assertHeader(t, req, lmw.CityHeader, "Munich")
	assertHeader(t, req, lmw.IPAddressHeader, ValidIP)

	req = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = fmt.Sprintf("%s:9999", ValidIPNoCity)
	instance.ServeHTTP(httptest.NewRecorder(), req)
	assertHeader(t, req, lmw.CountryCodeHeader, "US")
	assertHeader(t, req, lmw.RegionCodeHeader, lmw.Unknown)
	assertHeader(t, req, lmw.CityHeader, lmw.Unknown)
	assertHeader(t, req, lmw.IPAddressHeader, ValidIPNoCity)

	req = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = "qwerty:9999"
	instance.ServeHTTP(httptest.NewRecorder(), req)
	assertHeader(t, req, lmw.CountryCodeHeader, lmw.Unknown)
	assertHeader(t, req, lmw.RegionCodeHeader, lmw.Unknown)
	assertHeader(t, req, lmw.CityHeader, lmw.Unknown)
	assertHeader(t, req, lmw.IPAddressHeader, "qwerty")
}

func TestGeoIPFromXForwardedFor(t *testing.T) {
	mwCfg := mw.CreateConfig()
	mwCfg.CityDBPath = "data/tmp/GeoLite2-City.mmdb"
	mwCfg.PreferXForwardedForHeader = true

	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	mw.ResetLookup()
	instance, _ := mw.New(context.TODO(), next, mwCfg, "traefik-geoip")

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = fmt.Sprintf("%s:9999", ValidIP)
	req.Header.Set("X-Forwarded-For", ValidAlternateIP)
	instance.ServeHTTP(httptest.NewRecorder(), req)
	assertHeader(t, req, lmw.CountryCodeHeader, "DE")
	assertHeader(t, req, lmw.RegionCodeHeader, "BY")
	assertHeader(t, req, lmw.CityHeader, "Munich")
	assertHeader(t, req, lmw.IPAddressHeader, ValidAlternateIP)

	req = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = fmt.Sprintf("%s:9999", ValidIP)
	req.Header.Set("X-Forwarded-For", ValidAlternateIP+",188.193.88.100")
	instance.ServeHTTP(httptest.NewRecorder(), req)
	assertHeader(t, req, lmw.CountryCodeHeader, "DE")
	assertHeader(t, req, lmw.CountryHeader, "Germany")
	assertHeader(t, req, lmw.RegionCodeHeader, "BY")
	assertHeader(t, req, lmw.CityHeader, "Munich")
	assertHeader(t, req, lmw.IPAddressHeader, ValidAlternateIP)

	req = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = ValidIP + ":9999"
	req.Header.Set("X-Forwarded-For", "qwerty")
	instance.ServeHTTP(httptest.NewRecorder(), req)
	assertHeader(t, req, lmw.CountryHeader, lmw.Unknown)
	assertHeader(t, req, lmw.RegionHeader, lmw.Unknown)
	assertHeader(t, req, lmw.CityHeader, lmw.Unknown)
	assertHeader(t, req, lmw.IPAddressHeader, "qwerty")
}

func TestGeoIPCountryDBFromRemoteAddr(t *testing.T) {
	mwCfg := mw.CreateConfig()
	mwCfg.CountryDBPath = "data/tmp/GeoLite2-Country.mmdb"

	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	mw.ResetLookup()
	instance, _ := mw.New(context.TODO(), next, mwCfg, "traefik-geoip")

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	req.RemoteAddr = fmt.Sprintf("%s:9999", ValidIP)
	instance.ServeHTTP(httptest.NewRecorder(), req)

	assertHeader(t, req, lmw.CountryCodeHeader, "DE")
	assertHeader(t, req, lmw.CountryHeader, "Germany")
	assertHeader(t, req, lmw.IPAddressHeader, ValidIP)
}

// func TestGeoIpCityWithSpecialCharacters(t *testing.T) {
// 	mwCfg := mw.CreateConfig()
// 	mwCfg.CityDBPath = "data/tmp/GeoLite2-City.mmdb"

// 	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
// 	mw.ResetLookup()
// 	instance, _ := mw.New(context.TODO(), next, mwCfg, "traefik-geoip")

// 	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
// 	req.RemoteAddr = fmt.Sprintf("%s:9999", "179.96.134.192")
// 	instance.ServeHTTP(httptest.NewRecorder(), req)

// 	assertHeader(t, req, lmw.CountryCodeHeader, "BR")
// 	assertHeader(t, req, lmw.CountryHeader, "Brazil")
// 	assertHeader(t, req, lmw.CityHeader, "Marília")
// 	city := req.Header.Get(lmw.CityHeader)
// 	if utf8.ValidString(city) {
// 		log.Printf("City is valid UTF-8 %s", city)
// 	} else {
// 		log.Printf("City is NOT valid UTF-8 %s", city)
// 	}
// 	assertHeader(t, req, lmw.IPAddressHeader, "179.96.134.192")
// }

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()
	if req.Header.Get(key) != expected {
		t.Fatalf("invalid value of header [%s] != %s", key, req.Header.Get(key))
	}
}
