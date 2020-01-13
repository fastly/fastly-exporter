# fastly-exporter [![Latest Release](https://img.shields.io/github/release/peterbourgon/fastly-exporter.svg?style=flat-square)](https://github.com/peterbourgon/fastly-exporter/releases/latest) [![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fpeterbourgon%2Ffastly-exporter%2Fbadge&style=flat-square&label=build)](https://github.com/peterbourgon/fastly-exporter/actions?query=workflow%3ATest) [![Docker Status](https://img.shields.io/docker/build/mrnetops/fastly-exporter.svg)](https://hub.docker.com/r/mrnetops/fastly-exporter)

This program consumes from the [Fastly Real-time Analytics API][rt] and makes
the data available to [Prometheus][prom]. It should behave like you expect:
dynamically adding new services, removing old services, and reflecting changes
to service metadata like name and version.

[rt]: https://docs.fastly.com/api/analytics
[prom]: https://prometheus.io

## Getting

### Binary

Go to the [releases page][releases].

[releases]: https://github.com/peterbourgon/fastly-exporter/releases

### Docker

Avaliable as [mrnetops/fastly-exporter][container] from [Docker Hub][hub].

[container]: https://hub.docker.com/r/mrnetops/fastly-exporter
[hub]: https://hub.docker.com

```
docker pull mrnetops/fastly-exporter
```

### Source

If you have a working Go installation, you can install the latest revision from HEAD.

```
go get github.com/peterbourgon/fastly-exporter/cmd/fastly-exporter
```

## Using

### Basic

For simple use cases, all you need is a Fastly API token.
[See this link][token] for information on creating API tokens. 
The token can be provided via the `-token` flag or the 
`FASTLY_API_TOKEN` environment variable.

[token]: https://docs.fastly.com/guides/account-management-and-security/using-api-tokens#creating-api-tokens

```
fastly-exporter -token XXX
```

This will collect real-time stats for all Fastly services visible to your
token, and make them available as Prometheus metrics on 127.0.0.1:8080/metrics.

### Advanced

```
USAGE
  fastly-exporter [flags]

FLAGS
  -api-refresh 1m0s                        how often to poll api.fastly.com for updated service metadata
  -api-timeout 15s                         HTTP client timeout for api.fastly.com requests (5–60s)
  -debug false                             Log debug information
  -endpoint http://127.0.0.1:8080/metrics  Prometheus /metrics endpoint
  -name-exclude-regex ...                  if set, ignore any service whose name matches this regex
  -name-include-regex ...                  if set, only include services whose names match this regex
  -namespace fastly                        Prometheus namespace
  -rt-timeout 45s                          HTTP client timeout for rt.fastly.com requests (45–120s)
  -service ...                             if set, only include this service ID (repeatable)
  -shard ...                               if set, only include services whose hashed IDs modulo m equal n-1 (format 'n/m')
  -subsystem rt                            Prometheus subsystem
  -token ...                               Fastly API token (required; also via FASTLY_API_TOKEN)
  -version false                           print version information and exit
```

By default, all services available to your token will be exported. You can
specify an explicit set of service IDs by using the `-service xxx` flag.
(Service IDs are available at the top of your [Fastly dashboard][db].) You can
also include only those services whose name matches a regex by using the
`-name-include-regex '^Production'` flag, or reject any service whose name
matches a regex by using the `-name-exclude-regex '.*TEST.*'` flag.

[db]: https://manage.fastly.com/services/all

For tokens with access to a lot of services, it's possible to "shard" the
services among different instances of the fastly-exporter by using the `-shard`
flag. For example, to shard all services between 3 exporters, you would start
each exporter as

```
fastly-exporter [common flags] -shard 1/3
fastly-exporter [common flags] -shard 2/3
fastly-exporter [common flags] -shard 3/3
```

Flags which restrict the services that are exported combine with AND semantics.
That is, `-service A -service B -name-include-regex 'Foo'` would only export
data for service A and/or B if their names also matched "Foo". Or, specifying
`-name-include-regex 'Prod' -name-exclude-regex '^test-'` would only export data
for services whose names contained "Prod" and did not start with "test-".

### Docker

This repo contains a Dockerfile if you want to build and package it yourself.
You can also use a third-party Docker image.

```
docker run -p 8080:8080 mrnetops/fastly-exporter -token $MY_TOKEN
```

This repo also contains a [Docker Compose][compose] file, which boots up a full
fastly-exporter + [Prometheus][prom] + [Grafana][grafana] + Fastly dashboard
stack.

[compose]: https://github.com/docker/compose
[grafana]: https://grafana.com

```
env FASTLY_API_TOKEN=$MY_TOKEN docker-compose up
```

Access the Grafana dashboard at http://localhost:3000.

![Fastly Dashboard in Grafana](https://raw.githubusercontent.com/peterbourgon/fastly-exporter/master/compose/Fastly-Dashboard.png)
