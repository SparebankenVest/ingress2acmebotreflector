# Changelog

## Unreleased

- bump `golang.org/x/crypto` to `v0.52.0` to address CVE-2026-39830
- bump `golang.org/x/net` to `v0.55.0` to address CVE-2026-39821

## 0.1.16

- restructure `image` value into `repository`/`tag`/`pullPolicy`; image tag now
  defaults to the chart `appVersion`, removing version drift between the chart
  and the deployed image

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
