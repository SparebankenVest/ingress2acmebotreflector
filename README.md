# Helm

The helm repo for the helm chart `chart/ingress2acmebotreflector` is hosted by the github pages for this github repo. Add the repo locally:

```bash
helm repo add ingress2acmebotreflector-repo https://sparebankenvest.github.io/ingress2acmebotreflector/
```

## Helm chart installation

The definition of the helm chart is found in the `chart/ingress2acmebotreflector` directory.

### Install locally

1. Clone repo locally: `git clone git@github.com:SparebankenVest/ingress2acmebotreflector.git`
2. Go to chart directory: `cd ingress2acmebotreflector/chart`
3. Insert required values in `values.yaml`.
3. Install chart on cluster via helm: `helm install my-ingress2acmebotreflector ingress2acmebotreflector`

### Install using self hosted helm repo

1. Add helm repo: `helm repo add <REPO-NAME> <REPO-URL>`
2. Install helm chart on cluster: `helm install my-ingress2acmebotreflector <REPO-NAME>/ingress2acmebotreflector`
