## knctl route list

List routes

### Synopsis

List all routes in a namespace

```
knctl route list [flags]
```

### Examples

```

  # List all routes in namespace 'ns1'
  knctl route list -n ns1
```

### Options

```
  -h, --help               help for list
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
```

### Options inherited from parent commands

```
      --column strings              Filter to show only given columns
      --json                        Output as JSON
      --kubeconfig string           Path to the kubeconfig file ($KNCTL_KUBECONFIG or $KUBECONFIG)
      --kubeconfig-context string   Kubeconfig context override ($KNCTL_KUBECONFIG_CONTEXT)
      --no-color                    Disable colorized output
      --non-interactive             Don't ask for user input
      --tty                         Force TTY-like output
```

### SEE ALSO

* [knctl route](knctl_route.md)	 - Route management (create, curl, delete, list, show)

