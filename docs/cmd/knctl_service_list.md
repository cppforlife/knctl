## knctl service list

List services

### Synopsis

List all services in a namespace

```
knctl service list [flags]
```

### Examples

```

  # List all services in namespace 'ns1'
  knctl service list -n ns1
```

### Options

```
  -h, --help               help for list
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
```

### Options inherited from parent commands

```
      --column strings      Filter to show only given columns
      --json                Output as JSON
      --kubeconfig string   Path to the kubeconfig file ($KNCTL_KUBECONFIG or $KUBECONFIG)
      --no-color            Disable colorized output
      --non-interactive     Don't ask for user input
      --tty                 Force TTY-like output
```

### SEE ALSO

* [knctl service](knctl_service.md)	 - Service (annotate, delete, list, open)

