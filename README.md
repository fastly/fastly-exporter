# fastly-exporter [![Latest Release](https://img.shields.io/github/release/peterbourgon/fastly-exporter.svg?style=flat-square)](https://github.com/peterbourgon/fastly-exporter/releases/latest) [![builds.sr.ht status](https://builds.sr.ht/~peterbourgon/fastly-exporter.svg)](https://builds.sr.ht/~peterbourgon/fastly-exporter?) [![Docker Status](https://img.shields.io/docker/build/mrnetops/fastly-exporter.svg)](https://hub.docker.com/r/mrnetops/fastly-exporter)

This program consumes from the [Fastly Real-time Analytics API][rt] and makes
the data available to [Prometheus][prom]. It can provide metrics for every
service accessible to your API token, or an explicitly-specified set of
services. And it reflects when new services are created, old services are
deleted, or existing services have their names or versions updated.

[rt]: https://docs.fastly.com/api/analytics
[prom]: https://prometheus.io

## Getting

### Docker

Avaliable as [mrnetops/fastly-exporter][container] from [Docker Hub][hub].

[container]: https://hub.docker.com/r/mrnetops/fastly-exporter
[hub]: https://hub.docker.com

```
docker pull mrnetops/fastly-exporter
```

### Binary

Go to the [releases page][releases].

[releases]: https://github.com/peterbourgon/fastly-exporter/releases

### Source

If you have a working Go installation, you can install the latest revision from HEAD.

```
go get github.com/peterbourgon/fastly-exporter
```

## Using

```
USAGE
  fastly-exporter [flags]

FLAGS
  -debug false                             Log debug information
  -endpoint http://127.0.0.1:8080/metrics  Prometheus /metrics endpoint
  -namespace fastly                        Prometheus namespace
  -service ...                             Specific Fastly service ID (optional, repeatable)
  -subsystem rt                            Prometheus subsystem
  -token ...                               Fastly API token (required; also via FASTLY_API_TOKEN)
  -version false                           print version information and exit
```

A valid Fastly API token is mandatory. [See this link][token] for information
on creating API tokens. The token can also be provided via the `-token` flag
or the FASTLY_API_TOKEN environment variable.

Optional `-service` IDs can be specified to limit monitoring to specific
services. Service IDs are available at the top of your [Fastly dashboard][db].

[token]: https://docs.fastly.com/guides/account-management-and-security/using-api-tokens#creating-api-tokens
[db]: https://manage.fastly.com/services/all

### Docker

```
docker run -p 8080:8080 mrnetops/fastly-exporter -token $FASTLY_API_TOKEN
```

### Docker Compose

[Docker Compose][compose] for a full fastly-exporter + [Prometheus][prom] + 
[Grafana][grafana] + Fastly dashboard stack

[compose]: https://github.com/docker/compose
[grafana]: https://grafana.com

```
FASTLY_API_TOKEN=${FASTLY_API_TOKEN} docker-compose up
```

Access the [Grafana][grafana] dashboard via http://localhost:3000.

![Fastly Dashboard in Grafana](https://raw.githubusercontent.com/peterbourgon/fastly-exporter/master/compose/Fastly-Dashboard.png)
