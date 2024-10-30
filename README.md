# Traefik plugin for GeoIP

[Traefik](https://doc.traefik.io/traefik/) plugin 
that registers a custom middleware 
for getting data from 
[MaxMind GeoIP databases](https://www.maxmind.com/en/geoip2-services-and-databases) 
and pass it downstream via HTTP request headers.

Supports both 
[GeoIP2](https://www.maxmind.com/en/geoip2-databases) 
and 
[GeoLite2](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) databases.

## Installation 


### Docker Compose

Simple example to run in docker compose with auto download maxmind database

docker-compose.yaml
```yaml
services:

  traefik:
    image: "traefik:v3.2"
    container_name: "traefik"
    restart: unless-stopped
    command:
      #- "--log.level=DEBUG"
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entryPoints.web.address=:80"
      - "--experimental.localPlugins.traefikgeoip.moduleName=github.com/thiagotognoli/traefikgeoip"
    ports:
      - "80:80"
      - "8080:8080"
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
      - "../.:/plugins-local/src/github.com/thiagotognoli/traefikgeoip"
      - geoipupdate_data:/usr/share/GeoIP
    networks:
      - traefikgeoip_example

  whoami:
    image: "traefik/whoami"
    container_name: "traefik-whoami"
    restart: unless-stopped
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.whoami.rule=Host(`localhost`)"
      - "traefik.http.routers.whoami.entrypoints=web"
      - "traefik.http.middlewares.traefikgeoip.plugin.traefikgeoip.cityDbPath=/usr/share/GeoIP/GeoLite2-City.mmdb"
      - "traefik.http.middlewares.traefikgeoip.plugin.traefikgeoip.asnDbPath=/usr/share/GeoIP/GeoLite2-ASN.mmdb"
      - "traefik.http.middlewares.traefikgeoip.plugin.traefikgeoip.lightMode=true"
      - "traefik.http.middlewares.traefikgeoip.plugin.traefikgeoip.ipHeader=X-IP"
      - "traefik.http.middlewares.traefikgeoip.plugin.traefikgeoip.iso88591=true"
      # - "traefik.http.middlewares.traefikgeoip.plugin.traefikgeoip.preferXForwardedForHeader=true"
      # - "traefik.http.middlewares.traefikgeoip.plugin.traefikgeoip.failInError=true"
      - "traefik.http.routers.whoami.middlewares=traefikgeoip"
    networks:
      - traefikgeoip_example

  geoipupdate:
    container_name: cii_geoipupdate
    image: ghcr.io/maxmind/geoipupdate
    restart: unless-stopped
    # restart: on-failure
    env_file:
      - ./maxmind.env
    volumes:
      - geoipupdate_data:/usr/share/GeoIP
    networks:
      - traefikgeoip_example

volumes:
  geoipupdate_data:

networks:
  traefikgeoip_example:
    name: traefikgeoip_example
    driver: bridge
```

Request your free maxmind credentials in site and up in
- [MAXMIND](https://www.maxmind.com)
- [MAXMIND signup](https://www.maxmind.com/en/geolite2/signup?utm_source=kb&utm_medium=kb-link&utm_campaign=kb-create-account)

maxmind.env
```env
GEOIPUPDATE_ACCOUNT_ID=6025842
GEOIPUPDATE_LICENSE_KEY=abQ2rY_UMmru2Gd6LrbfGcrmKFwvR1duEDPZ_nmw
GEOIPUPDATE_EDITION_IDS=GeoLite2-ASN GeoLite2-City GeoLite2-Country
GEOIPUPDATE_FREQUENCY=72
```

### Kubernetes

The tricky part of installing this plugin into containerized environments, like Kubernetes,
is that a container should contain a database within it.


> [!WARNING]
> Setup below is provided for demonstration purpose and should not be used on production.
> Traefik's plugin site is observed to be frequently unavailable, 
> so plugin download may fail on pod restart.

Tested with [official Traefik chart](https://artifacthub.io/packages/helm/traefik/traefik) version 26.0.0.

The following snippet should be added to `values.yaml`:

```yaml
experimental:
  plugins:
    geoip2:
      moduleName: github.com/thiagotognoli/traefikgeoip
      version: v0.22.0
deployment:
  additionalVolumes:
    - name: geoip2
      emptyDir: {}
  initContainers:
    - name: download
      image: alpine
      volumeMounts:
        - name: geoip2
          mountPath: /tmp/geoip2
      command:
        - "/bin/sh"
        - "-ce"
        - |
          wget -P /tmp https://raw.githubusercontent.com/thiagotognoli/traefikgeoip/main/geolite2.tgz
          tar --directory /tmp/geoip2 -xvzf /tmp/geolite2.tgz
additionalVolumeMounts:
  - name: geoip2
    mountPath: /geoip2
```

### Create Traefik Middleware

```yaml
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: geoip2
  namespace: traefik
spec:
  plugin:
    geoip2:
      dbPath: "/geoip2/GeoLite2-City.mmdb"
```

## Configuration

The plugin currently supports the following configuration settings:

Name | Description
---- | ----
cityDbPath | Container path to City GeoIP database.
countryDbPath | Container path to Country GeoIP database.
asnDbPath | Container path to ASN GeoIP database.
preferXForwardedForHeader | Should `X-Forwarded-For` header be used to extract IP address. Default `false`.
ipHeader | Alternate Header of IP. Default `""`.
failInError | Not start plugin in error. Default `false`.
debug | Debug messages: false. Default `false`.
iso88591 | Encode in ISO-8859-1; Default: `false`.


## Development

Install Go, golangci-lint, yaegi and just

```sh
brew install go golangci-lint just
go install github.com/traefik/yaegi/cmd/yaegi@latest
```

To run linter and tests execute this command

```sh
just test
```

# Original Plugin

This project is a fork of https://github.com/traefik-plugins/traefikgeoip2

and some parts of https://github.com/Maronato/traefik_geoip
