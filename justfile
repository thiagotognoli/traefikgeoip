# list available receipes
@default:
  just --list

data:
  #!/usr/bin/env bash
  cd data
  mkdir -p mmdb
  go run main.go -i GeoLite2-City.json -o mmdb/GeoLite2-City.mmdb -t GeoLite2-City
  go run main.go -i GeoLite2-ASN.json -o mmdb/GeoLite2-ASN.mmdb -t GeoLite2-ASN
  go run main.go -i GeoLite2-Country.json -o mmdb/GeoLite2-Country.mmdb -t GeoLite2-Country

dist:
  #!/usr/bin/env bash
  ./utils/dist.sh

format:
  gofumpt -w -extra .

# lint go files
lint:
  golangci-lint run -v

# run regular golang tests
test-go:
  go test -v -cover ./...

@_clean-yaegi:
  rm -rf /tmp/yaegi*

# run tests via yaegi
test-yaegi: && _clean-yaegi
  #!/usr/bin/env bash
  set -euox

  TMP=$(mktemp -d yaegi.XXXXXX -p /tmp)
  WRK="${TMP}/go/src/github.com/thiagotognoli"
  mkdir -p ${WRK}
  ln -s `pwd` "${WRK}"
  cd "${WRK}/$(basename `pwd`)"
  env GOPATH="${TMP}/go" yaegi test -v .

# lint and test
test: format lint test-go test-yaegi

clean:
  rm -rf *.mmdb
  rm -rf data/tmp

vendor:
  go mod vendor

