# fastly-exporter [![Latest Release](https://img.shields.io/github/release/peterbourgon/fastly-exporter.svg?style=flat-square)](https://github.com/peterbourgon/fastly-exporter/releases/latest) [![builds.sr.ht status](https://builds.sr.ht/~peterbourgon/fastly-exporter.svg)](https://builds.sr.ht/~peterbourgon/fastly-exporter?) [![Docker Status](https://img.shields.io/docker/build/mrnetops/fastly-exporter.svg)](https://hub.docker.com/r/mrnetops/fastly-exporter)

This program consumes from the [Fastly Real-time Analytics API][rt] and makes
the data available to [Prometheus][prom]. It should behave like you expect:
dynamically adding new services, removing old services, and reflecting changes
to service metadata like name and version.

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
  -api-refresh 1m0s                        how often to poll api.fastly.com for updated service metadata
  -debug false                             Log debug information
  -endpoint http://127.0.0.1:8080/metrics  Prometheus /metrics endpoint
  -name-regex ...                          if provided, only export services whose names match this regex
  -namespace fastly                        Prometheus namespace
  -service ...                             if set, only export services with this service ID (repeatable)
  -shard ...                               if set, only export services whose hashed IDs mod m equal (n-1) (format 'n/m')
  -subsystem rt                            Prometheus subsystem
  -token ...                               Fastly API token (required; also via FASTLY_API_TOKEN)
  -version false                           print version information and exit
```

A valid Fastly API token is mandatory. [See this link][token] for information on
creating API tokens. The token can be provided via the `-token` flag or the
FASTLY_API_TOKEN environment variable.

[token]: https://docs.fastly.com/guides/account-management-and-security/using-api-tokens#creating-api-tokens
[db]: https://manage.fastly.com/services/all

There are a number of ways to determine which services are exported. By default,
all services available to your token will be exported. You can specify an
explicit set of service IDs by using the `-service xxx` flag. (Service IDs are
available at the top of your [Fastly dashboard][db].) You can also specify only
those services whose user-defined name matches a regex by using the 
`-name-regex '^Production'` flag.

For tokens with access to a lot of services, it's possible to "shard" the
services among different instances of the fastly-exporter by using the `-shard`
flag. For example, to shard all services between 3 expoters, you would start
each exporter as

```
fastly-exporter ... -shard 1/3
fastly-exporter ... -shard 2/3
fastly-exporter ... -shard 3/3
```

Flags which restrict the services that are exported combine with AND semantics.
That is, `-service A -service B -name-regex 'Foo'` would only export data for
service A or service B if their names also matched 'Foo'.

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
