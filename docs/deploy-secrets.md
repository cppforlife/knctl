## Deploy with secrets

See [Basic Workflow](./basic-workflow.md) for introduction.

Below tutorial shows how to reference Kubernetes secrets and/or config maps from your Knative service.

Create new namespace

```bash
$ kubectl create ns deploy-with-secrets

$ export KNCTL_NAMESPACE=deploy-with-secrets
```

```bash
$ kubectl -n deploy-with-secrets create secret generic simple-msg --from-literal=val=123
```

Deploy service that builds image from a Git repo, and then deploys it

```bash
$ knctl deploy \
    --service simple-app \
    --image gcr.io/knative-samples/helloworld-go \
    --env-secret SIMPLE_MSG=simple-msg/val
```

(Similar can be done via `--env-config-map` flag to reference config map values.)

Curl the deployed service and see that it responds with `123`

```bash
$ knctl curl --service simple-app

Running: curl '-H' 'Host: simple-app.default.example.com' 'http://0.0.0.0:80'

Hello World: 123!
```
