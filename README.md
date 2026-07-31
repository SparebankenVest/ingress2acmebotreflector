# ingress2acmebotreflector

A simple k8s controller that works with the <https://github.com/shibayan/keyvault-acmebot> project to move cert orders from inside the cluster to the outside.
The application works from inside the cluster by watching all ingress resources and every time an ingress resource is changed or created it checks the acmebot api if a
tls cert exists for the given ingress, and if not it orders a cert for the domain in question.

The image for the controller is published to GitHub Container Registry: <https://github.com/SparebankenVest/ingress2acmebotreflector/pkgs/container/ingress2acmebotreflector>

## Description

// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

1. Install [keyvault acmebot](https://github.com/shibayan/keyvault-acmebot) in your azure environment
2. Set up a the keyvault acmebot the expose an api through the azure portal
3. Grant the managed-identity of the k8s cluster where the controller is access to the acmebot api
4. Install the controller in your `k8s cluster`` by using either the helm chart or the sample deployment file in the project.

### Build image

1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/deployment/samples/
```

2. Build your image:

```sh
make docker-build docker-push IMG=<some-registry>/ingress2acmebotreflector:tag
```
3. Push your image to the image repository location specified by `IMG`:

```sh
make docker-push docker-push IMG=<some-registry>/ingress2acmebotreflector:tag
```

### Deploy to cluster from terminal

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/ingress2acmebotreflector:tag
```


### Uninstall CRDs

To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller

UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Helm
The helm repo for the helm chart `chart/ingress2acmebotreflector` is hosted by the github pages for this github repo. Add the repo locally:
```
helm repo add ingress2acmebotreflector-repo https://sparebankenvest.github.io/ingress2acmebotreflector/
```

### Helm chart installation
The definition of the helm chart is found in the `chart/ingress2acmebotreflector` directory. 
#### Install locally
1. Clone repo locally: `git clone git@github.com:SparebankenVest/ingress2acmebotreflector.git`
2. Go to chart directory: `cd ingress2acmebotreflector/chart`
3. Insert required values in `values.yaml`. 
3. Install chart on cluster via helm: `helm install my-ingress2acmebotreflector ingress2acmebotreflector`
#### Install using self hosted helm repo
1. Add helm repo: `helm repo add <REPO-NAME> <REPO-URL>`
2. Install helm chart on cluster: `helm install my-ingress2acmebotreflector <REPO-NAME>/ingress2acmebotreflector`
## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out

1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Build and Release

Both the application image and the Helm chart are versioned from a single
source of truth: the `appVersion` field in
`config/charts/ingress2acmebotreflector/Chart.yaml`. The chart deploys the image
tag matching `appVersion` by default (see `image.tag` in `values.yaml`), so you
do not maintain the image tag separately.

### Versioning model

- **App image version** = `appVersion` in `Chart.yaml`.
- **Chart version** = `version` in `Chart.yaml` (bump on any chart change).

There are no manual `git tag` steps — tags and GitHub Releases are created
automatically when a version is merged to `main`.

### Release flow

1. Create a branch and bump the version(s) in
   `config/charts/ingress2acmebotreflector/Chart.yaml`:
   - Bump `appVersion` when the application code changes.
   - Bump `version` for any chart change (required for a new chart release).
   - Add a matching entry to the chart `CHANGELOG.md`.
2. Open a PR against `main`. On the PR:
   - `docker-ci.yaml` builds and pushes a **pre-release image** tagged
     `ghcr.io/sparebankenvest/ingress2acmebotreflector:<appVersion>-pre.<PR#>.<shortsha>`
     that you can pull and test. The exact tag is shown in the workflow summary.
   - `helm-ci.yaml` lints and install-tests the chart, and guards against image
     tag drift (`image.tag` must stay empty so it tracks `appVersion`).
3. Merge to `main`. On merge:
   - `docker-publish.yaml` reads `appVersion`. If that image tag does not already
     exist, it builds a **multi-arch** (`linux/amd64,linux/arm64`) image, pushes
     `ghcr.io/sparebankenvest/ingress2acmebotreflector:<appVersion>`, creates the
     git tag `<appVersion>`, and cuts a GitHub Release with generated notes. If
     the tag already exists (a merge that did not bump `appVersion`), this step is
     skipped.
   - `helm-release.yaml` packages and publishes the chart via
     `chart-releaser` when `version` in `Chart.yaml` was bumped (unchanged
     versions are skipped).

### Pre-release / release-candidate images

Every PR push produces a unique immutable `-pre.<PR#>.<shortsha>` image for
testing. These are only built for PRs; the final semver image is published only
on merge to `main`.

### Workflows summary

| Workflow | Trigger | Purpose |
| --- | --- | --- |
| `docker-ci.yaml` | PR to `main` | Build + push `<appVersion>-pre.<PR#>.<sha>` image |
| `docker-publish.yaml` | Push to `main` | Publish `<appVersion>` image, git tag, GitHub Release |
| `helm-ci.yaml` | PR to `main` | Lint/test chart + image-tag drift guard |
| `helm-release.yaml` | Push to `main` | Package + publish Helm chart |

