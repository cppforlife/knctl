## Standalone Build

See [Basic Workflow](./basic-workflow.md) for introduction.

Create new namespace

```bash
$ knctl create namespace -n standalone-build

$ export KNCTL_NAMESPACE=standalone-build
```

Create Docker Hub secret where images will be pushed (assumes that pushing needs authentication, but pulling does not)

```bash
$ knctl create basic-auth-secret -s docker-reg1 --docker-hub -u <your-username> -p <your-password>
```

Create service account that references above credential

```bash
$ knctl create service-account -a serv-acct1 -s docker-reg1
```

Deploy service that builds image from a Git repo, and then deploys it

```bash
$ knctl build \
    --build build1 \
    --git-url https://github.com/cppforlife/simple-app \
    --git-revision master \
    --service-account serv-acct1 \
    --image index.docker.io/<your-username>/<your-repo>
```
