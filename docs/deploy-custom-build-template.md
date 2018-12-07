## Deploy with custom Build Template (for example Buildpack)

See [Basic Workflow](./basic-workflow.md) for introduction.

Create new namespace

```bash
$ kubectl create ns deploy-with-buildpack

$ export KNCTL_NAMESPACE=deploy-with-buildpack
```

Install Buildpack Build Template

```bash
$ kubectl -n deploy-with-buildpack apply -f https://raw.githubusercontent.com/knative/build-templates/master/buildpack/buildpack.yaml
```

(See [Buildpack build template](https://github.com/knative/build-templates/tree/master/buildpack) for more info.)

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
$ knctl service-account create -a serv-acct1 -s docker-reg1 [-s docker-reg2]
```

Deploy service that builds image from a Git repo, and then deploys it

```bash
$ knctl deploy \
    --service simple-app \
    --git-url https://github.com/cppforlife/simple-app \
    --git-revision master \
    --template buildpack \
    --template-env GOPACKAGENAME=main \
    --service-account serv-acct1 \
    --image index.docker.io/<your-username>/<your-repo> \
    --env SIMPLE_MSG=123
```
