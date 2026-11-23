# Deprecated and unmaintained

> **Warning**
> This repository is no longer maintained and has been archived. For upstream Prometheus benchmarking infrastructure, see [prometheus/test-infra](https://github.com/prometheus/test-infra).

---

## This repository contains tools and configuration files for the testing and benchmarking used in the Prometheus project.

### [`/prombench`](/prombench)

Prombench is a project for automated E2E testing and benchmarking for Prometheus.

See [prombench/README.md](prombench/README.md) for full description.

## Building tools from source
With a working go modules enabled Go environment:
- Install [promu](https://github.com/prometheus/promu): `go install https://github.com/prometheus/promu`
- `promu build`
