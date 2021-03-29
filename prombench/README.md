# Automated Prometheus E2E testing and benchmarking.

![Prombench Design](design.svg)

It runs with [Github Actions](https://github.com/features/actions) on a [Google Kubernetes Engine Cluster](https://cloud.google.com/kubernetes-engine/).
It is designed to support adding more k8s providers.

## Overview of the manifest files

The `/manifest` directory contains all the kubernetes manifest files.
- `cluster_gke.yaml` : This is used to create the Main Node in gke.
- `cluster_eks.yaml` : This is used to create the Main Node in eks.
- `cluster-infra/` : These are the persistent components of the Main Node.
- `prombench/` : These resources are created and destroyed for each prombench test.

## Setup and run prombench

Prombench can be run on various providers, following are the provider specific instructions:
    
- Instructions for [Google Kubernetes Engine](docs/gke.md)
- Instructions for [Kubernetes In Docker](docs/kind.md)
- Instructions for [Elastic Kubernetes Service](docs/eks.md)

## Setup GitHub Actions

Place a workflow file in the `.github` directory of the repository.
See the [prometheus/prometheus](https://github.com/prometheus/prometheus) repository for an example.

Create a github action `TEST_INFRA_PROVIDER_AUTH` secret with the base64 encoded content of the `AUTH_FILE`.

```
cat $AUTH_FILE | base64 -w 0
```

### Trigger tests via a Github comment.
<!-- If you change the heading, also change the anchor in the comment monitor config map. -->

---

> Due to the high cost of each test, only maintainers can manage tests.

**Starting:**

- `/prombench main` or `/prombench master` - compare PR with the main/master branch.
- `/prombench v2.4.0` - compare PR with a release version, from [quay.io/prometheus/prometheus:releaseVersion](https://quay.io/prometheus/prometheus:releaseVersion)

**Restarting:**

- `/prombench restart <release_version>`

**Stopping:**

- `/prombench cancel`

### Building Docker Image

```
docker build -t prominfra/prombench:master .
```
