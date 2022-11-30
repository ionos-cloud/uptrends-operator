# :chart_with_upwards_trend: uptrends Operator

[![Release](https://github.com/ionos-cloud/uptrends-operator/actions/workflows/release.yml/badge.svg)](https://github.com/ionos-cloud//uptrends-operator/actions/workflows/release.yml)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)

This operator helps to configure [uptrends](https://www.uptrends.com/) monitoring for your [Ingress Controllers](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/).

## Introduction

The operator is based on the [uptrends](https://github.com/ionos-cloud/uptrends-go) package. It is a Kubernetes operator that watches for Ingress resources and creates uptrends checks for them. It also watches for changes in the Ingress resources and updates the uptrends checks accordingly.

## Helm

[Helm](https://helm.sh/) can be used to install :chart_with_upwards_trend: uptrends Operator.

```bash
helm repo add uptrends https://ionos-cloud.github.io/uptrends-operator/
helm repo update
```

The most recent version is installed via.

```bash
helm install uptrends uptrends/uptrends --create-namespace --namespace uptrends --version v0.0.3
```

The required `API_USERNAME` and `API_PASSWORD` can be securely configured via `envFrom` in the `values.yaml`.

## Environment

### `API_USERNAME` 

This configures the required username for the uptrends API access. See the [uptrends](https://www.uptrends.com/support/kb/api) documentation for more information.

### `API_PASSWORD` 

This configures the required password for the uptrends API access. See the [uptrends](https://www.uptrends.com/support/kb/api) documentation for more information.

## Annotations

The operator supports creating a monitor via the `Uptrends` kind, but also via annotations on an `Ingress`. The following keys are supported.

###  `uptrends.ionos-cloud.github.io/monitor.type` Default: `HTTPS`

This can be either `HTTPS` or `HTTP`.

### `uptrends.ionos-cloud.github.io/monitor.interval` Default: `"5"`

This can be an interval from `1` to `60` minutes.

> The annotations are evaluates agains the `host` fields on the `rules`. Wildcard hosts and empty hosts are ignored.

### `uptrends.ionos-cloud.github.io/monitor.guid` Default: `""`

This can be used to add the monitor to a Monitor Group identified with a `MonitorGroupID`.

### `uptrends.ionos-cloud.github.io/monitor.regions` Default:`""`

This is a list of regions to include as checkpoints. An example `"54,1007"`

### `uptrends.ionos-cloud.github.io/monitor.checkpoints` Default:`""`

This is a list of point of presence to include.

### `uptrends.ionos-cloud.github.io/monitor.exclude` Default: `""`

This is a list of point of presence to exclude as checkpoints.

## Examples

[/examples](/examples/) contains the example of an `Ingress` and `Uptrends` based monitor.

> Renaming the ingress does not delete monitors.

## License

[Apache 2.0](/LICENSE)
