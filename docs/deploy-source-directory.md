## Deploy from local source directory

See [Basic Workflow](./basic-workflow.md) for introduction.

Grab example source code via Git

```bash
$ git clone https://github.com/cppforlife/simple-app

$ cd simple-app
```

Create new namespace

```bash
$ knctl namespace create -n deploy-from-source

$ export KNCTL_NAMESPACE=deploy-from-source
```

Create Docker Hub secret for pushing images

```bash
$ knctl basic-auth-secret create -s docker-reg1 --docker-hub -u <your-username> -p <your-password>
```

If necessary, create Docker Hub secret for pulling images

```bash
$ knctl basic-auth-secret create -s docker-reg2 --docker-hub -u <your-username> -p <your-password> --for-pulling
```

Create service account that references above credentials

```bash
$ knctl service-account create -a serv-acct1 -s git1 -s docker-reg1 [-s docker-reg2]
```

Deploy service that builds image from a Git repo, and then deploys it

```bash
$ knctl deploy \
    --service simple-app \
    --directory=$PWD \
    --service-account serv-acct1 \
    --image index.docker.io/<your-username>/<your-repo> \
    --env SIMPLE_MSG=123
```

Optionally make some changes (in `./app.go` for example)

Deploy with updated source without commiting this change to Git

```bash
$ knctl deploy \
    --service simple-app \
    --directory=$PWD \
    --service-account serv-acct1 \
    --image index.docker.io/<your-username>/<your-repo> \
    --env SIMPLE_MSG=123
```
