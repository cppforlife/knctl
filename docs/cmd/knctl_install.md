## knctl install

Install Knative and Istio

### Synopsis

Install Knative and Istio.

Requires 'kubectl' command installed on a the system.

```
knctl install [flags]
```

### Options

```
  -m, --exclude-monitoring   Exclude installation of monitoring components
  -h, --help                 help for install
  -p, --node-ports           Use service type NodePorts instead of type LoadBalancer
```

### Options inherited from parent commands

```
      --column strings      Filter to show only given columns
      --json                Output as JSON
      --kubeconfig string   Path to the kubeconfig file ($KNCTL_KUBECONFIG) (default "/Users/pivotal/.kube/config")
      --no-color            Disable colorized output
      --non-interactive     Don't ask for user input
      --tty                 Force TTY-like output
```

### SEE ALSO

* [knctl](knctl.md)	 - knctl controls Knative resources (annotate, build, create, curl, delete, deploy, install, list, logs, open, route, tag, uninstall, untag, version)

