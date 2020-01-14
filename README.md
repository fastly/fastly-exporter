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

### Filtering services

By default, all services available to your token will be exported. You can
specify an explicit set of service IDs to export by using the `-service xxx`
flag. (Service IDs are available at the top of your [Fastly dashboard][db].) You
can also include only those services whose name matches a regex by using the
`-service-whitelist '^Production'` flag, or skip any service whose name matches
a regex by using the `-service-blacklist '.*TEST.*'` flag.

[db]: https://manage.fastly.com/services/all

For tokens with access to a lot of services, it's possible to "shard" the
services among different instances of the fastly-exporter by using the
`-service-shard` flag. For example, to shard all services between 3 exporters,
you would start each exporter as

```
fastly-exporter [common flags] -service-shard 1/3
fastly-exporter [common flags] -service-shard 2/3
fastly-exporter [common flags] -service-shard 3/3
```

### Filtering exported metrics

By default, all metrics provided by the Fastly real-time stats API are exported
as Prometheus metrics. You can export only those metrics whose name matches a
regex by using the `-metric-whitelist bytes_total$` flag, or elide any metric
whose name matches a regex by using the `-metric-blacklist imgopto` flag.

### Filter semantics

All flags that restrict services or metrics are repeatable. Repeating the same
flag causes its condition to be combined with OR semantics. For example,
`-service A -service B` would include both services A and B (but not service C).
Or, `-service-blacklist Test -service-blacklist Staging` would skip any service
whose name contained Test or Staging.

Different flags (for the same filter target) combine with AND semantics. For
example, `-metric-whitelist 'bytes_total$' -metric-blacklist imgopto` would only
export metrics whose names ended in bytes_total, but didn't include imgopto.

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
