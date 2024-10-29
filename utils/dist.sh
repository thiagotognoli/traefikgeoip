#!/usr/bin/env bash
scriptPath="$(sourcePath=`readlink -f "$0"` && echo "${sourcePath%/*}")"
basePath="$(cd $scriptPath/.. && pwd)"

cd $basePath

[[ -e "$basePath/dist" ]] && rm -rf "$basePath/dist"


path="dist/plugins-local/src/github.com/thiagotognoli/traefikgeoip"

mkdir -p "$path"
cp go.mod go.sum .traefik.yml middleware.go "$path"
cp -R geoip2 lib "$path"

