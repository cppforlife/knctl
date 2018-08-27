## knctl service open

Open web browser pointing at a service domain

### Synopsis

Open web browser pointing at a service domain.

Requires 'open' command installed on the system.

```
knctl service open [flags]
```

### Examples

```

  # Open web browser pointing at service 'svc1' in namespace 'ns1'
  knctl service open -s svc1 -n ns1
```

### Options

```
  -h, --help               help for open
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -s, --service string     Specified service
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

