# fastly-exporter [![Latest Release](https://img.shields.io/github/release/fastly/fastly-exporter.svg?style=flat-square)](https://github.com/fastly/fastly-exporter/releases/latest) [![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Ffastly%2Ffastly-exporter%2Fbadge%3Fref%3Dmain&style=flat-square)](https://actions-badge.atrox.dev/fastly/fastly-exporter/goto?ref=main)

This program consumes from the [the Fastly API][api] and makes
the data available to [Prometheus][prom]. It should behave like you expect:
dynamically adding new services, removing old services, and reflecting changes
to service metadata like name and version.

Fastly APIs consumed:

* [Real-time Analytics API][rt]
* [Service list][svc]
* [Custom TLS Certificates][certs]
* [POPs][pops]

[api]: https://www.fastly.com/documentation/reference/api/
[rt]: https://www.fastly.com/documentation/reference/api/metrics-stats/realtime/
[svc]: https://www.fastly.com/documentation/reference/api/services/service/#list-services
[pops]: https://www.fastly.com/documentation/reference/api/utils/pops/
[certs]: https://www.fastly.com/documentation/reference/api/tls/custom-certs/

[prom]: https://prometheus.io

## Getting

### Binary

Go to the [releases page][releases].

[releases]: https://github.com/fastly/fastly-exporter/releases

### Docker

Available on the [packages page][pkg] as [fastly/fastly-exporter][img].

[pkg]: https://github.com/fastly/fastly-exporter/packages
[img]: https://github.com/fastly/fastly-exporter/pkgs/container/fastly-exporter

```sh
docker pull ghcr.io/fastly/fastly-exporter:latest
```

Note that version `latest` will track RCs, alphas, etc. -- always use an
explicit version in production.

### Helm chart

[Helm](https://helm.sh) must be installed to use the [prometheus-community/fastly-exporter](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus-fastly-exporter) chart.
Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

Once Helm is set up properly, add the repo as follows:

```sh
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```

And install:

```sh
helm upgrade --install fastly-exporter prometheus-fastly-exporter --namespace monitoring --set token="fastly_api_token"
```

### Source

If you have a working Go installation, you can clone the repo and install the
binary from any revision, including HEAD.

```sh
git clone git@github.com:fastly/fastly-exporter
cd fastly-exporter
go build ./cmd/fastly-exporter
./fastly-exporter -h
```

## Using

### Authentication

A valid Fastly API token is required to use the exporter. [See this link][token]
for information on creating API tokens. The token can be provided via the
`-token` flag or the `FASTLY_API_TOKEN` environment variable.

If you would like to export TLS certificate metrics the token will need TLS Management permissions:
* If you are using a [user token][token] the user must have TLS Management permissions. [See this link][user_perms] for how to update this access permissions.
* If you are using an [automation token][auto_token] it must have TLS Management permissions.

[token]: https://www.fastly.com/documentation/guides/account-info/account-management/using-api-tokens/#creating-api-tokens
[auto_token]: https://www.fastly.com/documentation/guides/account-info/account-management/using-api-tokens/#creating-automation-tokens
[user_perms]: https://www.fastly.com/documentation/guides/account-info/user-access-and-control/configuring-user-roles-and-permissions/#changing-user-roles-and-access-permissions-for-existing-users

### Basic

For simple use cases, all you need is a Fastly API token.

```sh
fastly-exporter -token XXX
```

This will collect real-time stats for all Fastly services visible to your token,
and make them available as Prometheus metrics on [127.0.0.1:8080/metrics][local].

[local]: http://127.0.0.1:8080/metrics

### Filtering services

By default, all services available to your token will be exported. You can
specify an explicit set of service IDs to export by using the `-service xxx`
flag. (Service IDs are available at the top of your [Fastly dashboard][db].) You
can also include only those services whose name matches a regex by using the
`-service-allowlist '^Production'` flag, or exclude any service whose name matches
a regex by using the `-service-blocklist '.*TEST.*'` flag.

[db]: https://manage.fastly.com/services/all

For tokens with access to a lot of services, it's possible to "shard" the
services among different fastly-exporter instances by using the `-service-shard`
flag. For example, to shard all services between 3 exporters, you would start
each exporter as

```sh
fastly-exporter [common flags] -service-shard 1/3
fastly-exporter [common flags] -service-shard 2/3
fastly-exporter [common flags] -service-shard 3/3
```

### Filtering metrics

By default, all metrics provided by the Fastly real-time stats API are exported
as Prometheus metrics. You can export only those metrics whose name matches a
regex by using the `-metric-allowlist 'bytes_total$'` flag, or exclude any metric
whose name matches a regex by using the `-metric-blocklist imgopto` flag.

### Filter semantics

All flags that filter services or metrics are repeatable. Repeating the same
flag causes its condition to be combined with OR semantics. For example,
`-service A -service B` would include both services A and B (but not service C).
Or, `-service-blocklist Test -service-blocklist Staging` would skip any service
whose name contained Test or Staging.

Different flags (for the same filter target) combine with AND semantics. For
example, `-metric-allowlist 'bytes_total$' -metric-blocklist imgopto` would only
export metrics whose names ended in bytes_total, but didn't include imgopto.

### Metrics Grouping: by datacenter or aggregate

The Fastly real-time stats API returns measurements grouped by datacenter as
well as aggregated measurements for all datacenters. By default, exported
metrics are grouped by datacenter. The response body size of the metrics
endpoint can potentially be very large. This will be exacerbated when using
the exporter with many services, many origins with Origin Inspector, and many
domains with Domain Inspector. One way to reduce the output size of the
metrics endpoint is by using the `-aggregate-only` flag. When this flag is
used only the `aggregated` metrics from the real-time stats API will be
exported. Metrics will still include the datacenter label but it will always
be set to "aggregate".

### Service discovery

Per-service metrics are available via `/metrics?target=<service ID>`. Available
services are enumerated as targets on the `/sd` endpoint, which is compatible
with the [generic HTTP service discovery][httpsd] feature of Prometheus. An
example Prometheus scrape config for the Fastly exporter follows.

[httpsd]: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#http_sd_config

```yaml
scrape_configs:
  - job_name: fastly-exporter
    http_sd_configs:
      - url: http://127.0.0.1:8080/sd
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: service
      - target_label: __address__
        replacement: 127.0.0.1:8080
```

### Dashboards and Alerting

Data from the the Fastly exporter can be used to build dashboards and alerts with [Grafana][grafana] and [Alertmanager][alertmanager]. For a fully working example see [fastly-dashboards][dashboards] created by [@mrnetops][mrnetops]. Fastly-dashboards contains a Docker Compose setup, which boots up a full fastly-exporter + Prometheus + Alertmanager + Grafana + Fastly dashboard stack with Slack alerting integration.

[grafana]: https://grafana.com
[alertmanager]: https://prometheus.io/docs/alerting/latest/alertmanager/
[mrnetops]: https://github.com/mrnetops
[dashboards]: https://github.com/mrnetops/fastly-dashboards
