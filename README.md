# fastly-exporter [![Latest Release](https://img.shields.io/github/release/peterbourgon/fastly-exporter.svg?style=flat-square)](https://github.com/peterbourgon/fastly-exporter/releases/latest)

This program consumes from the [Fastly Real-time Analytics API][rt] and makes
the data available to [Prometheus][prom].

[rt]: https://docs.fastly.com/api/analytics
[prom]: https://prometheus.io

## Getting

Go to the [releases page](/releases), or, if you have a working Go installation,
you can install the latest revision from HEAD.

```
go get github.com/peterbourgon/fastly-exporter
```

## Using

```
USAGE
  fastly-exporter [flags]

FLAGS
  -debug false                             log debug information
  -endpoint http://127.0.0.1:8080/metrics  Prometheus /metrics endpoint
  -namespace ...                           Prometheus namespace
  -service ...                             Fastly service ID (repeatable)
  -subsystem ...                           Prometheus subsystem
  -token ...                               Fastly API token

VERSION
  1.0.0
```

A valid API -token and at least one -service ID are mandatory. Your service ID
is available at the top of your [Fastly dashboard][db]. [See this link][token]
for information on creating API tokens.

[db]: https://manage.fastly.com/services/all
[token]: https://docs.fastly.com/guides/account-management-and-security/using-api-tokens#creating-api-tokens
