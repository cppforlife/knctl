## knctl list builds

List builds

### Synopsis

List all builds in a namespace

```
knctl list builds [flags]
```

### Examples

```

  # List all builds in namespace 'ns1'
  knctl list builds -n ns1
```

### Options

```
  -h, --help               help for builds
  -n, --namespace string   Specified namespace (or default from kubeconfig)
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

