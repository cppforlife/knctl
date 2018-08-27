# Knative CLI

Knative CLI (`knctl`) provides simple set of commands to interact with a [Knative installation](https://github.com/knative/docs).

Grab pre-built binaries from the [Releases page](https://github.com/cppforlife/knctl/releases).

## Docs

- [Basic workflow](./docs/basic-workflow.md)
- [Deploy from public Git repo](./docs/deploy-public-git-repo.md)
- [Deploy from private Git repo](./docs/deploy-private-git-repo.md)
- [Deploy from local source directory](./docs/deploy-source-directory.md)
- [Deploy with custom Build Template (for example Buildpack)](./docs/deploy-custom-build-template.md)
- [`knctl` as a `kubectl` plugin](./docs/kubectl-plugin.md)
- Advanced
  - [Manage domains](./docs/manage-domains.md)
  - [Traffic splitting WIP](./docs/traffic-splitting.md)
  - [Stanalone build](./docs/standalone-build.md)
  - [Annotations](./docs/annotations.md)
  - [Ingresses](./docs/ingresses.md)
  - [Complete command reference](./docs/cmd/knctl.md)

## Development

If you are interested in contributing, see [CONTRIBUTING.md](./CONTRIBUTING.md).

```bash
# export GOPATH=...

$ ./hack/build.sh

$ ./knctl version
```
