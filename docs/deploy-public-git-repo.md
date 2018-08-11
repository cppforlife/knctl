## Deploy from public Git repo

See [Basic Workflow](./basic-workflow.md) for introduction.

Create new namespace

```bash
$ knctl create namespace -n deploy-from-git

$ export KNCTL_NAMESPACE=deploy-from-git
```

Create Docker Hub secret where images will be pushed

```bash
$ knctl create basic-auth-secret -s docker-reg1 --docker-hub -u <your-username> -p <your-password>
```

If Docker repository is private, there is no need to create this secret (and use below); otherwise, create Docker Hub secret for pulling images

```bash
$ kubectl -n deploy-from-git create secret docker-registry docker-reg2 \
    --docker-server https://index.docker.io \
    --docker-username <your-username> \
    --docker-password <your-password> \
    --docker-email noop
```

Create service account that references above credentials

```bash
$ knctl create service-account -a serv-acct1 -s docker-reg1 [ -p docker-reg2 ]
```

Deploy service that builds image from a Git repo, and then deploys it

```bash
$ knctl deploy \
    --service simple-app \
    --git-url https://github.com/cppforlife/simple-app \
    --git-revision master \
    --service-account serv-acct1 \
    --image index.docker.io/<your-username>/<your-repo> \
    --env SIMPLE_MSG=123
```
