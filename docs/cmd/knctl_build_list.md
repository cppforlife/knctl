## knctl build list

List builds

### Synopsis

List all builds in a namespace

```
knctl build list [flags]
```

### Examples

```

  # List all builds in namespace 'ns1'
  knctl build list -n ns1
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

* [knctl build](knctl_build.md)	 - Build management (create, delete, list, show)

