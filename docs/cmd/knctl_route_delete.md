## knctl route delete

Delete route

### Synopsis

Delete route

```
knctl route delete [flags]
```

### Examples

```

  # Delete route 'route1' in namespace 'ns1'
  knctl route delete --route route1 -n ns1
```

### Options

```
  -h, --help               help for delete
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
      --route string       Specified route
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

* [knctl route](knctl_route.md)	 - Route management (create, delete, list)

