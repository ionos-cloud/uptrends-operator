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

## License

[Apache 2.0](/LICENSE)
