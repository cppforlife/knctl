## Manage Domains

See [Basic Workflow](./basic-workflow.md) for introduction.

List available domains

```bash
$ knctl list domains

Domains

Name                  Default
my-other-domain.test  true

1 domains

Succeeded
```

Change default domain

```bash
$ knctl create domain my-domain.test --default
```

Deploy sample service with new default domain

```bash
$ knctl -n default deploy --service hello --image gcr.io/knative-samples/helloworld-go --env TARGET=123
```
