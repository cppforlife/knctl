# Knative CLI

Knative CLI (`knctl`) provides simple set of commands to interact with a [Knative installation](https://github.com/knative/docs).

Grab pre-built binaries from the [Releases page](https://github.com/cppforlife/knctl/releases).

## Build

```bash
# export GOPATH=...

$ ./hack/build.sh

$ ./knctl version
```

## Docs

- [Basic workflow](./docs/basic-workflow.md)
- [Deploy from public Git repo](./docs/deploy-public-git-repo.md)
- [Traffic splitting](./docs/traffic-splitting.md)
- [Misc](./docs/misc.md)
- [Complete command reference](./docs/cmd/knctl.md).

If you are interested in contributing, see [CONTRIBUTING.md](./CONTRIBUTING.md).
