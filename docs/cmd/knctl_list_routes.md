## knctl list routes

List routes

### Synopsis

List all routes in a namespace

```
knctl list routes [flags]
```

### Examples

```

  # List all routes in namespace 'ns1'
  knctl list routes -n ns1
```

### Options

```
  -h, --help               help for routes
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE)
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

* [knctl list](knctl_list.md)	 - List resources (services, revisions, builds, pods, ingresses)

