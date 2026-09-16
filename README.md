# fastly-exporter [![Latest Release](https://img.shields.io/github/release/fastly/fastly-exporter.svg?style=flat-square)](https://github.com/fastly/fastly-exporter/releases/latest) [![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Ffastly%2Ffastly-exporter%2Fbadge%3Fref%3Dmain&style=flat-square)](https://actions-badge.atrox.dev/fastly/fastly-exporter/goto?ref=main)

This program consumes from the [the Fastly API][api] and makes
the data available to [Prometheus][prom]. It should behave like you expect:
dynamically adding new services, removing old services, and reflecting changes
to service metadata like name and version.

Fastly APIs consumed:

* [Real-time Analytics API][rt]
* [Origin Inspector Real-time API][oi-rt] (when the account is entitled, see [below](#origin-inspector-and-domain-inspector))
* [Domain Inspector Real-time API][di-rt] (when the account is entitled)
* Product entitlement (`GET /entitled-products/{product}`), to decide whether to poll Origin Inspector and Domain Inspector
* [Service list][svc]
* [Custom TLS Certificates][certs]
* [POPs][pops]

[api]: https://www.fastly.com/documentation/reference/api/
[rt]: https://www.fastly.com/documentation/reference/api/metrics-stats/realtime/
[oi-rt]: https://www.fastly.com/documentation/reference/api/metrics-stats/origin-inspector/real-time/
[di-rt]: https://www.fastly.com/documentation/reference/api/metrics-stats/domain-inspector/real-time/
[svc]: https://www.fastly.com/documentation/reference/api/services/service/#list-services
[pops]: https://www.fastly.com/documentation/reference/api/utils/pops/
[certs]: https://www.fastly.com/documentation/reference/api/tls/custom-certs/

[prom]: https://prometheus.io

# Installation

## Binary

Go to the [releases page][releases].

[releases]: https://github.com/fastly/fastly-exporter/releases

## Docker

Available on the [packages page][pkg] as [fastly/fastly-exporter][img].

[pkg]: https://github.com/fastly/fastly-exporter/packages
[img]: https://github.com/fastly/fastly-exporter/pkgs/container/fastly-exporter

```sh
docker pull ghcr.io/fastly/fastly-exporter:latest
```

Note that version `latest` will track RCs, alphas, etc. -- always use an
explicit version in production.

## Helm chart

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

## Source

If you have a working Go installation, you can clone the repo and install the
binary from any revision, including HEAD.

```sh
git clone git@github.com:fastly/fastly-exporter
cd fastly-exporter
go build ./cmd/fastly-exporter
./fastly-exporter -h
```

# Using the Exporter

## Authentication

A valid Fastly API token is required to use the exporter. [See this link][token]
for information on creating API tokens. The token can be provided via the
`-token` flag or the `FASTLY_API_TOKEN` environment variable.

If you would like to export TLS certificate metrics the token will need TLS Management permissions:
* If you are using a [user token][token] the user must have TLS Management permissions. [See this link][user_perms] for how to update this access permissions.
* If you are using an [automation token][auto_token] it must have TLS Management permissions.

[token]: https://www.fastly.com/documentation/guides/account-info/account-management/using-api-tokens/#creating-api-tokens
[auto_token]: https://www.fastly.com/documentation/guides/account-info/account-management/using-api-tokens/#creating-automation-tokens
[user_perms]: https://www.fastly.com/documentation/guides/account-info/user-access-and-control/configuring-user-roles-and-permissions/#changing-user-roles-and-access-permissions-for-existing-users

## Basic

For simple use cases, all you need is a Fastly API token.

```sh
fastly-exporter -token XXX
```

This will collect real-time stats for all Fastly services visible to your token,
and make them available as Prometheus metrics on [127.0.0.1:8080/metrics][local].

[local]: http://127.0.0.1:8080/metrics

## Origin Inspector and Domain Inspector

If your account is entitled to [Origin Inspector][oi] or [Domain Inspector][di],
the exporter polls their real-time APIs alongside the standard real-time stats
and exports the results. There is no flag to turn this on. At startup, and every
`-product-refresh` (default 10m), the exporter calls
`GET https://api.fastly.com/entitled-products/{product}` for each product and
starts one extra subscriber per service for each product the account has access
to. The log shows this as `type=origin_inspector subscriber=create` and
`type=domain_inspector subscriber=create`. Run with `-debug` to see the
entitlement result itself (`product=origin_inspector hasAccess=true`).

Both products are [enabled per service][oi-enable]. Services without the product
enabled don't produce origin or domain metrics.

[oi]: https://docs.fastly.com/products/origin-inspector
[di]: https://docs.fastly.com/products/domain-inspector
[oi-enable]: https://www.fastly.com/documentation/reference/api/products/origin_inspector/

### Origin Inspector metrics

Exported under the `origin` subsystem (`fastly_origin_*` with the default
namespace). Every metric carries the labels `service_id`, `service_name`,
`datacenter`, `origin`, and `source`. `source` is one of `delivery`, `compute`,
or `waf`, identifying which part of the Fastly platform made the origin request.

| Metric | Type | Extra labels |
|--------|------|--------------|
| `fastly_origin_responses_total` | counter | |
| `fastly_origin_resp_body_bytes_total` | counter | |
| `fastly_origin_resp_header_bytes_total` | counter | |
| `fastly_origin_status_code_total` | counter | `status_code` |
| `fastly_origin_status_group_total` | counter | `status_group` |
| `fastly_origin_latency_seconds` | histogram | |

### Domain Inspector metrics

Exported under the `domain` subsystem (`fastly_domain_*`). Every metric carries
the labels `service_id`, `service_name`, `datacenter`, and `domain`.

| Metric | Type | Extra labels |
|--------|------|--------------|
| `fastly_domain_requests_total` | counter | |
| `fastly_domain_edge_requests_total` | counter | |
| `fastly_domain_edge_hit_requests_total` | counter | |
| `fastly_domain_edge_miss_requests_total` | counter | |
| `fastly_domain_edge_hit_ratio` | gauge | |
| `fastly_domain_resp_body_bytes_total` | counter | |
| `fastly_domain_resp_header_bytes_total` | counter | |
| `fastly_domain_edge_resp_body_bytes_total` | counter | |
| `fastly_domain_edge_resp_header_bytes_total` | counter | |
| `fastly_domain_bereq_body_bytes_total` | counter | |
| `fastly_domain_bereq_header_bytes_total` | counter | |
| `fastly_domain_origin_fetches` | counter | |
| `fastly_domain_origin_fetch_resp_body_bytes` | counter | |
| `fastly_domain_origin_fetch_resp_header_bytes` | counter | |
| `fastly_domain_origin_offload` | gauge | |
| `fastly_domain_status_code_total` | counter | `status_code` |
| `fastly_domain_status_group_total` | counter | `status_group` |
| `fastly_domain_origin_status_code_total` | counter | `status_code` |
| `fastly_domain_origin_status_group_total` | counter | `status_group` |
| `fastly_domain_http2_total` | counter | |
| `fastly_domain_http3_total` | counter | |
| `fastly_domain_tls_total` | counter | `tls_version` |

### Controlling cardinality

Origin metrics are emitted per datacenter, per origin, per source. Domain
metrics are emitted per datacenter, per domain. On a service with many origins
or domains, this multiplies the size of the `/metrics` response by the number of
Fastly POPs. Options, from least to most aggressive:

```sh
# Drop the per-datacenter breakdown and keep only aggregated values
fastly-exporter -token XXX -aggregate-only

# Keep origin status codes, drop the other origin metrics
fastly-exporter -token XXX -metric-blocklist '^fastly_origin_(resp_|responses|latency)'

# Turn domain metrics off entirely
fastly-exporter -token XXX -metric-blocklist '^fastly_domain_'
```

Blocklisting a metric stops it from being exported. It doesn't stop the
exporter from polling rt.fastly.com for that product.

## Filtering services

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

## Filtering metrics

By default, all metrics provided by the Fastly real-time stats API are exported
as Prometheus metrics. You can export only those metrics whose name matches a
regex by using the `-metric-allowlist 'bytes_total$'` flag, or exclude any metric
whose name matches a regex by using the `-metric-blocklist imgopto` flag.

## Filter semantics

All flags that filter services or metrics are repeatable. Repeating the same
flag causes its condition to be combined with OR semantics. For example,
`-service A -service B` would include both services A and B (but not service C).
Or, `-service-blocklist Test -service-blocklist Staging` would skip any service
whose name contained Test or Staging.

Different flags (for the same filter target) combine with AND semantics. For
example, `-metric-allowlist 'bytes_total$' -metric-blocklist imgopto` would only
export metrics whose names ended in bytes_total, but didn't include imgopto.

## Metrics Grouping: by datacenter or aggregate

The Fastly real-time stats API returns measurements grouped by datacenter as
well as aggregated measurements for all datacenters. By default, exported
metrics are grouped by datacenter. The response body size of the metrics
endpoint can potentially be very large. This will be exacerbated when using
the exporter with many services, many origins with Origin Inspector, and many
domains with Domain Inspector (see [Controlling cardinality](#controlling-cardinality)).
One way to reduce the output size of the
metrics endpoint is by using the `-aggregate-only` flag. When this flag is
used only the `aggregated` metrics from the real-time stats API will be
exported. Metrics will still include the datacenter label but it will always
be set to "aggregate".

## Service discovery

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

## Dashboards and Alerting

Data from the the Fastly exporter can be used to build dashboards and alerts with [Grafana][grafana] and [Alertmanager][alertmanager]. For a fully working example see [fastly-dashboards][dashboards] created by [@mrnetops][mrnetops]. Fastly-dashboards contains a Docker Compose setup, which boots up a full fastly-exporter + Prometheus + Alertmanager + Grafana + Fastly dashboard stack with Slack alerting integration.

[grafana]: https://grafana.com
[alertmanager]: https://prometheus.io/docs/alerting/latest/alertmanager/
[mrnetops]: https://github.com/mrnetops
[dashboards]: https://github.com/mrnetops/fastly-dashboards

# Maintainers

## Releases

To make a new release:

1. Push a new git tag, e.g. `git tag v9.4.0 -m'v9.4.0'; git push --tags`
1. GitHub Actions drafts the new release
1. Click the "[Generate release notes](https://docs.github.com/en/repositories/releasing-projects-on-github/automatically-generated-release-notes#creating-automatically-generated-release-notes-for-a-new-release)" in GitHub
1. Publish the new release
