## Deploy from public Git repo

See [Basic Workflow](./basic-workflow.md) for introduction.

Create new namespace

```bash
$ knctl create namespace -n deploy-from-git

$ export KNCTL_NAMESPACE=deploy-from-git
```

Create Docker Hub secret for pushing images

```bash
$ knctl create basic-auth-secret -s docker-reg1 --docker-hub -u <your-username> -p <your-password>
```

If necessary, create Docker Hub secret for pulling images

```bash
$ knctl create basic-auth-secret -s docker-reg2 --docker-hub -u <your-username> -p <your-password> --for-pulling
```

Create service account that references above credentials

```bash
$ knctl create service-account -a serv-acct1 -s docker-reg1 [-s docker-reg2]
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

See how to [deploy from private Git repo](./deploy-private-git-repo.md).
