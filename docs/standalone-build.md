## Standalone Build

See [Basic Workflow](./basic-workflow.md) for introduction.

Create new namespace

```bash
$ kubectl create ns standalone-build

$ export KNCTL_NAMESPACE=standalone-build
```

Create Docker Hub secret for pushing images

```bash
$ knctl basic-auth-secret create -s docker-reg1 --docker-hub -u <your-username> -p <your-password>
```

Create service account that references above credential

```bash
$ knctl service-account create -a serv-acct1 -s docker-reg1
```

Deploy service that builds image from a Git repo, and then deploys it

```bash
$ knctl build create \
    --build build1 \
    --git-url https://github.com/cppforlife/simple-app \
    --git-revision master \
    --service-account serv-acct1 \
    --image index.docker.io/<your-username>/<your-repo>
```
