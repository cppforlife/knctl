## Traffic Splitting WIP

**Currently does not work with latest Knative Serving**

### Tags

Each tag identifies single revision within a service. CLI uses Kubernetes labels to store tag information. Tags can be used to reference particular revision when changing traffic configuration. By default, CLI will apply two tags: `latest` and `previous`.

```bash
$ knctl -n default tag revision --revision hello:latest -t stable

$ knctl -n default list revisions --service hello
```

### Routing

```bash
$ export KNCTL_NAMESPACE=default

# Deploy a new revision of sample service without changing traffic configuration
$ knctl deploy --service hello --image gcr.io/knative-samples/helloworld-go --env TARGET=123

# Direct 10% of traffic to new revision
$ knctl route --route hello -p hello:latest=10% -p hello:previous=90%

# Direct a bit more of traffic to new revision
$ knctl route --route hello -p hello:latest=40% -p hello:previous=60%

# Direct all traffic to latest revision
$ knctl route --route hello -p hello:latest=100%

# Director all traffic back to previous revision
$ knctl route --route hello -p hello:previous=100%
```
