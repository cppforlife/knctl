## knctl service delete

Delete service

### Synopsis

Delete service

```
knctl service delete [flags]
```

### Examples

```

  # Delete service 'svc1' in namespace 'ns1'
  knctl service delete -s svc1 -n ns1
```

### Options

```
  -h, --help               help for delete
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -s, --service string     Specified service
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

