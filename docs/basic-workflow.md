## Basic Workflow

Install `knctl` by grabbing pre-built binaries from the [Releases page](https://github.com/cppforlife/knctl/releases)

```bash
$ shasum -a 265 ~/Downloads/knctl-*
# Compare checksum output to what's included in the release notes

$ mv ~/Downloads/knctl-* /usr/local/bin/knctl

$ chmod +x /usr/local/bin/knctl
```

Install Istio and Knative

```bash
$ knctl install
```

Deploy sample service (to current namespace)

```bash
$ knctl deploy --service hello --image gcr.io/knative-samples/helloworld-go --env TARGET=123
```

List deployed services

```bash
$ knctl service list

Services in namespace 'default'

Name   Domain                     Annotations  Age
hello  hello.default.example.com  -            1d

1 services
```

Check that at least one Pod is in `Running` state

```bash
$ knctl pod list --service hello

Pods for service 'hello'

Revision     Name                                    Phase    Restarts  Age
hello-00001  hello-00001-deployment-c9cc8b88c-8hw4x  Running  0         10s

1 pods
```

Curl the deployed service and see that it responds

```bash
$ knctl curl --service hello

Running: curl '-H' 'Host: hello.default.example.com' 'http://0.0.0.0:80'

Hello World: 123!
```

Fetch recent logs from the deployed service

```bash
$ knctl logs -f --service hello
hello-00001 > hello-00001-deployment-7d4b4c5cc-v6jvl | 2018/08/02 17:21:51 Hello world sample started.
hello-00001 > hello-00001-deployment-7d4b4c5cc-v6jvl | 2018/08/02 17:22:04 Hello world received a request.
```

Change environment variable and see changes were applied

```bash
$ knctl deploy --service hello --image gcr.io/knative-samples/helloworld-go --env TARGET=new-value

$ knctl curl --service hello

Running: curl '-H' 'Host: hello.default.example.com' 'http://0.0.0.0:80'

Hello World: new-value!
```

List multiple revisions of the deployed service

```bash
$ knctl revision list --service hello

Revisions for service 'hello'

Name         Allocated Traffic %  Serving State  Age
hello-00002  100%                 Active         2m
hello-00001  0%                   Reserve        3m

2 revisions
```

See how to:

- [deploy from public Git repo](./deploy-public-git-repo.md)
- [deploy from local source directory](./deploy-source-directory.md)
