#!/usr/bin/env bash
scriptPath="$(sourcePath=`readlink -f "$0"` && echo "${sourcePath%/*}")"
basePath="$(cd $scriptPath/.. && pwd)"

cd $basePath

[[ -e "$basePath/dist" ]] && rm -rf "$basePath/dist"


path="dist/plugins-local/src/github.com/thiagotognoli/traefikgeoip"

mkdir -p "$path"
cp go.mod go.sum .traefik.yml middleware.go "$path"
mkdir -p "$path/lib"
cp lib/*.go "$path/lib"/.
mkdir -p "$path/geoip2"
cp geoip2/*.go "$path/geoip2"/.
mkdir -p "$path/geoip2_iso88591"
cp geoip2_iso88591/*.go "$path/geoip2_iso88591"/.

