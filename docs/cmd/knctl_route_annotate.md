## knctl route annotate

Annotate route

### Synopsis

Annotate route

```
knctl route annotate [flags]
```

### Examples

```

  # Annotate route 'rt1' in namespace 'ns1' with key and value
  knctl route annotate --route rt1 -a key=value -n ns1
```

### Options

```
  -a, --annotation strings   Set annotation (format: key=value) (can be specified multiple times)
  -h, --help                 help for annotate
  -n, --namespace string     Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
      --route string         Specified route
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

* [knctl route](knctl_route.md)	 - Route management (annotate, curl, delete, list, show)

