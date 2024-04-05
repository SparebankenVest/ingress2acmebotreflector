# ingress2acmebotreflector

A simple k8s controller that works with the <https://github.com/shibayan/keyvault-acmebot> project to move cert orders from inside the cluster to the outside.
The application works from inside the cluster by watching all ingress resources and every time an ingress resource is changed or created it checks the acmebot api if a
tls cert exists for the given ingress, and if not it orders a cert for the domain in question.

The image for the controller can be found on dockerhub: <https://hub.docker.com/r/spvest/ingress2acmebotreflector>

## Description

// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

1. Install [keyvault acmebot](https://github.com/shibayan/keyvault-acmebot) in your azure environment
2. Set up a the keyvault acmebot the expose an api through the azure portal
3. Grant the managed-identity of the k8s cluster where the controller is access to the acmebot api, this can be done using the script in the scripts folder.
4. Install the controller in your `k8s cluster`` by using the sample deployment file in the project.

### Build image

1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
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