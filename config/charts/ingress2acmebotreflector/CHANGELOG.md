# Changelog

## 0.1.24

### Security

- disable automatic service account token mounting in deployment spec (`automountServiceAccountToken: false`)

## 0.1.23

- bump `github.com/onsi/ginkgo/v2` to `v2.34.0`

## 0.1.22

- bump Azure SDK `azcore` to `v1.24.0`
- update the pinned `softprops/action-gh-release` digest to `8a0cd65`

## 0.1.21

- bump Kubernetes Go modules (`api`, `apimachinery`, and `client-go`) to `v0.37.2`
  and `sigs.k8s.io/controller-runtime` to `v0.25.2`
- raise the Go module requirement to `1.26.0` and update the Docker builder to Go `1.26`
- update the pinned `docker/login-action` digest to `8aa5831`

## 0.1.20

- bump Azure SDK `azcore` to `v1.23.3` and `azidentity` to `v1.15.0`
- update Helm to `v4.3.1`, `azure/setup-helm` to `v5.0.1`, and
  `helm/kind-action` to `v1.15.1`
- update Docker build workflows: `docker/build-push-action` to `v7` and
  `docker/login-action`, `docker/setup-buildx-action`, and
  `docker/setup-qemu-action` to `v4`

## 0.1.19

- update `actions/checkout` to `v7` and `softprops/action-gh-release` to `v3`

## 0.1.18

- update `helm/chart-testing-action` to `v2.9.0` and `actions/setup-python` to `v7`

## 0.1.17

- configure Renovate for dependency updates, remove Dependabot configuration,
  and add `.gitignore`
- pin the distroless runtime image to digest `25ac8f0`
- bump `github.com/onsi/ginkgo/v2` to `v2.33.1`,
  `github.com/onsi/gomega` to `v1.44.0`, and update transitive Go dependencies
- update Helm to `v3.22.1`, `azure/setup-helm` to `v4.3.1`, and
  `helm/chart-releaser-action` to `v1.9.0`

## 0.1.16

- restructure `image` value into `repository`/`tag`/`pullPolicy`; image tag now
  defaults to the chart `appVersion`, removing version drift between the chart
  and the deployed image

### Security Patch

- bump `golang.org/x/crypto` to `v0.52.0` to address CVE-2026-39830
- bump `golang.org/x/net` to `v0.55.0` to address CVE-2026-39821

## 0.1.15

- add default container security context
- add default CPU request

## 0.1.14

- add support for additional pod labels via `podLabels`
- align chart `appVersion` with release `0.1.14`

## 0.1.13

- change chcr.io repository. Doesn't support capital letters.

## 0.1.12

- add pipeline for pushing docker images to ghcr.io
- syncing helmchart with new version
