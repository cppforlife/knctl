# Knative CLI

Knative CLI (`knctl`) provides simple set of commands to interact with a Knative installation.

If you are interested in contributing, see [CONTRIBUTING.md](./CONTRIBUTING.md).

## Basic Workflow

Install Istio and Knative

```bash
$ knctl install
```

Deploy sample service

```bash
$ knctl -n default deploy --service hello --image gcr.io/knative-samples/helloworld-go --env TARGET=123
```

List deployed services

```bash
$ knctl -n default list services

Services in namespace 'default'

Name   Domain                     Internal Domain                  Created At
hello  hello.default.example.com  hello.default.svc.cluster.local  2018-08-01T17:35:51-07:00

1 services
```

Curl the deployed service and see that it responds

```bash
$ knctl -n default curl --service hello

Running: curl '-H' 'Host: hello.default.example.com' 'http://0.0.0.0:80'

Hello World: 123!
```

Fetch recent logs from the deployed service

```bash
$ knctl -n default logs -f --service hello
hello-00001 > hello-00001-deployment-7d4b4c5cc-v6jvl | 2018/08/02 17:21:51 Hello world sample started.
hello-00001 > hello-00001-deployment-7d4b4c5cc-v6jvl | 2018/08/02 17:22:04 Hello world received a request.
```

Change environment variable and see changes were applied

```bash
$ knctl -n default deploy --service hello --image gcr.io/knative-samples/helloworld-go --env TARGET=new-value

$ knctl -n default curl --service hello

Running: curl '-H' 'Host: hello.default.example.com' 'http://0.0.0.0:80'

Hello World: new-value!
```

List multiple revisions of the deployed service

```bash
$ knctl -n default list revisions

Revisions for service 'hello'

Name         Allocated Traffic %  Serving State  Created At
hello-00002  100%                 Active         2018-08-01T17:35:51-07:00
hello-00001  0%                   Reserve        2018-08-01T17:32:51-07:00

2 revisions
```

## Annotations

Annotations can be used to store helpful metadata about a revision.

```
$ export KNCTL_NAMESPACE=default

$ knctl annotate revision -a owner=cppforlife -a version=1.3.3 --revision latest

$ knctl list revisions --service hello
```

## Tags

Each tag identifies single revision within a service. CLI uses Kubernetes labels to store tag information. Tags can be used to reference particular revision when changing traffic configuration. By default, CLI will apply two tags: `latest` and `previous`.

```
$ knctl -n default tag revision --revision hello:latest -t stable

$ knctl list revisions --service hello
```

## Traffic Splitting

**Currently non-functional**

```
$ export KNCTL_NAMESPACE=default

# Deploy a new revision of sample service without changing traffic configuration
$ knctl deploy --service hello --image gcr.io/knative-samples/helloworld-go --env TARGET=123 --traffic-unchanged

# Direct 10% of traffic to new revision
$ knctl route --route hello -p hello:latest=10% -p hello:previous=90%

# Direct a bit more of traffic to new revision
$ knctl route --route hello -p hello:latest=40% -p hello:previous=60%

# Direct all traffic to latest revision
$ knctl route --route hello -p hello:latest=100%

# Director all traffic back to previous revision
$ knctl route --route hello -p hello:previous=100%
```

## Misc

```bash
$ knctl -n default list ingresses

Ingresses

Name                    Addresses  Ports  Created At
knative-ingressgateway  0.0.0.0    80     2018-08-01T17:10:38-07:00
                                   443
                                   32400

1 ingresses
```

## Command Reference

See [complete command reference](./docs/cmd/knctl.md).

## TODO

- rollout and rollback command
- top command
