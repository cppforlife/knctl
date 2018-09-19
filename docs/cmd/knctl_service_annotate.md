## knctl service annotate

Annotate service

### Synopsis

Annotate service

```
knctl service annotate [flags]
```

### Examples

```

  # Annotate service 'srv1' in namespace 'ns1' with key and value
  knctl service annotate -s srv1 -a key=value -n ns1
```

### Options

```
  -a, --annotation strings   Set annotation (format: key=value) (can be specified multiple times)
  -h, --help                 help for annotate
  -n, --namespace string     Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -s, --service string       Specified service
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

* [knctl service](knctl_service.md)	 - Service management (annotate, delete, list, open, show)

