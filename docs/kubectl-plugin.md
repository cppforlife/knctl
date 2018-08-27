## `knctl` as a `kubectl` plugin

`kubectl` recently [introduced new CLI plugin mechanism](https://github.com/kubernetes/kubernetes/pull/66876) on master branch (not officially released yet).

To try it out, build new `kubectl` locally from master

```bash
$ git clone https://github.com/kubernetes/kubernetes ~/workspace/kubernetes/src/k8s.io/kubernetes

$ export GOPATH=~/workspace/kubernetes/src

$ go build k8s.io/kubernetes/cmd/kubectl

$ ./kubectl version
```

Install `knctl` by grabbing pre-built binaries from the [Releases page](https://github.com/cppforlife/knctl/releases)

```bash
$ shasum -a 265 ~/Downloads/knctl-*
# Compare checksum output to what's included in the release notes

$ mv ~/Downloads/knctl-* /usr/local/bin/kubectl-kn

$ chmod +x /usr/local/bin/kubectl-kn
```

`kubectl` will find any binary named `kubectl-*` on your `PATH` and consider it as a plugin

```bash
$ ./kubectl plugin list

/usr/local/bin/kubectl-kn
```

List Knative services

```bash
$ ./kubectl kn service list

Services in namespace 'default'

Name   Domain                     Annotations  Age
hello  hello.default.example.com  -            1d

1 services
```
