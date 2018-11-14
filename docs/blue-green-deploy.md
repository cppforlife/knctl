## Blue-green deploy

Deploy first version of your service and lock down route to point to this revision.

```bash
$ knctl deploy -s hello -i gcr.io/knative-samples/helloworld-go -e TARGET=first --managed-route=false
$ knctl rollout --route hello -p hello:latest=100%
```

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
