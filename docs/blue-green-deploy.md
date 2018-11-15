## Blue-green deploy

Requires `knctl` v0.1.0+ and Knative 0.2.0+.

[Blue-green deployments](https://martinfowler.com/bliki/BlueGreenDeployment.html) can be performed by controlling Knative route objects with `knctl rollout` command.

Deploy first version of your service and lock down route to point to this revision.

```bash
$ knctl deploy -s hello -i gcr.io/knative-samples/helloworld-go -e TARGET=first --managed-route=false
$ knctl rollout --route hello -p hello:latest=100%
```

(Note `--managed-route=false` flag that indicates to knctl that Knative route will be controlled via subsequent commands.)

Deploy another version of your service with a newly developed feature.

```bash
$ knctl deploy -s hello -i gcr.io/knative-samples/helloworld-go -e TARGET=second --managed-route=false
```

At this point route still points to the old revision. Let's roll out new version to 10% of users.

```bash
$ knctl rollout --route hello -p hello:latest=10% -p hello:previous=90%
```

Once appropriate metrics are verified that new version is OK, roll out remaining traffic.

```bash
$ knctl rollout --route hello -p hello:latest=100%
```

### Tags

Each tag identifies single revision within a service. CLI uses Kubernetes labels to store tag information. Tags can be used to reference particular revision when changing traffic configuration. By default, CLI will apply two tags: `latest` and `previous`.

```bash
$ knctl revision tag --revision hello:latest -t stable
$ knctl revision list --service hello
```
