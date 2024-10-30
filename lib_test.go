package traefikgeoip_test

// func TestLibUtils(t *testing.T) {
// 	mwCfg := mw.CreateConfig()
// 	mwCfg.CityDBPath = "data/mmdb/GeoLite2-City.mmdb"
// 	// mwCfg.Iso88591 = true

// 	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
// 	mw.ResetLookup()
// 	instance, _ := mw.New(context.TODO(), next, mwCfg, "traefik-geoip")

// 	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
// 	req.RemoteAddr = fmt.Sprintf("%s:9999", "179.96.134.192")
// 	instance.ServeHTTP(httptest.NewRecorder(), req)

// 	city := req.Header.Get(lmw.CityHeader)
// 	if city != lmw.StringIso88591ToUtf8(lmw.StringUtf8ToIso88591(city)) {
// 		t.Fatalf("error converting string o iso and back to utf. string:  [%s]", city)
// 	}
// }
