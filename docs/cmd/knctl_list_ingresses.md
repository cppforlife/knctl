## knctl list ingresses

List ingresses

### Synopsis

List all ingresses labeled as `knative: ingressgateway` in Istio's namespace

```
knctl list ingresses [flags]
```

### Options

```
  -h, --help   help for ingresses
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

* [knctl list](knctl_list.md)	 - List (builds, domains, ingresses, pods, revisions, routes, services)

