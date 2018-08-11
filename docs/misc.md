## Annotations

Annotations can be used to store helpful metadata about a service or revision.

```
$ export KNCTL_NAMESPACE=default

$ knctl annotate service -a owner=cppforlife --revision hello

$ knctl annotate revision -a version=1.3.3 --revision hello:latest

$ knctl list revisions --service hello
```

## Ingresses

```bash
$ knctl -n default list ingresses

Ingresses

Name                    Addresses  Ports  Age
knative-ingressgateway  0.0.0.0    80     1d
                                   443
                                   32400

1 ingresses
```
