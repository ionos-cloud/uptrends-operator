# :chart_with_upwards_trend: uptrends Operator

[![Release](https://github.com/ionos-cloud/uptrends-operator/actions/workflows/release.yml/badge.svg)](https://github.com/ionos-cloud//uptrends-operator/actions/workflows/release.yml)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)

This operator helps to configure [uptrends](https://www.uptrends.com/) monitoring for your [Ingress Controllers](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/).

:warning: this is experimental work :test_tube: and interfaces may change.

## Introduction

The operator is based on the [uptrends](https://github.com/ionos-cloud/uptrends-go) package. It is a Kubernetes operator that watches for Ingress resources and creates uptrends checks for them. It also watches for changes in the Ingress resources and updates the uptrends checks accordingly.

## Environment

### `API_USERNAME` 

This configures the required username for the uptrends API access. See the [uptrends](https://www.uptrends.com/support/kb/api) documentation for more information.

### `API_PASSWORD` 

This configures the required password for the uptrends API access. See the [uptrends](https://www.uptrends.com/support/kb/api) documentation for more information.

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

## License

[Apache 2.0](/LICENSE)
