# Knative CLI

Knative CLI (`knctl`) provides simple set of commands to interact with a [Knative installation](https://github.com/knative/docs).

Grab pre-built binaries from the [Releases page](https://github.com/cppforlife/knctl/releases).

## Docs

- [Basic workflow](./docs/basic-workflow.md)
- [Deploy from public Git repo](./docs/deploy-public-git-repo.md)
- [Deploy from private Git repo](./docs/deploy-private-git-repo.md)
- [Deploy from local source directory](./docs/deploy-source-directory.md)
- [Deploy with custom Build Template (for example Buildpack)](./docs/deploy-custom-build-template.md)
- [Deploy with secrets](./docs/deploy-secrets.md)
- [Blue-green deploy](./docs/blue-green-deploy.md)
- [`knctl` as a `kubectl` plugin](./docs/kubectl-plugin.md)
- Advanced
  - [Manage domains](./docs/manage-domains.md)
  - [Standalone build](./docs/standalone-build.md)
  - [Annotations](./docs/annotations.md)
  - [Ingresses](./docs/ingresses.md)
  - [Complete command reference](./docs/cmd/knctl.md)
- Blog posts
  - [IBM Developer Blog: Introducing Knctl: A simpler way to work with Knative](https://developer.ibm.com/blogs/2018/11/12/knctl-a-simpler-way-to-work-with-knative/)
  - starkandwayne.com blog
	  - [Deploying 12-factor apps to Knative](https://www.starkandwayne.com/blog/deploying-12factor-apps-to-knative/)
	  - [Building and deploying applications to Knative](https://starkandwayne.com/blog/building-and-deploying-applications-to-knative/)
	  - [Adding public traffic to Knative on Google Kubernetes Engine](https://starkandwayne.com/blog/public-traffic-into-knative-on-gke/)
	  - [Adding a custom hostname domain for Knative services](https://starkandwayne.com/blog/adding-a-custom-domain-for-knative-services/)
	  - [Build Docker images inside your Kubernetes with Knative Build](https://starkandwayne.com/blog/build-docker-images-inside-kubernetes-with-knative-build/)
	  - [Binding secrets to Knative services](https://starkandwayne.com/blog/binding-secrets-to-knative-services/)
- Talks
  - [Introducing Knctl, a command line tool for Knative (YouTube)](https://www.youtube.com/watch?v=cJyJGm22Pf4)
  - [Dr Nic's Introducing Knative to Small Teams talk](https://speakerdeck.com/drnic/introducing-knative-to-small-teams) (slides only)

## Development

```bash
# export GOPATH=...

$ ./hack/build.sh

$ ./knctl version
```
