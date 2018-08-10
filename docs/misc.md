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

Name                    Addresses  Ports  Created At
knative-ingressgateway  0.0.0.0    80     2018-08-01T17:10:38-07:00
                                   443
                                   32400

1 ingresses
```
