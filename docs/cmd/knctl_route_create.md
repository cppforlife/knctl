## knctl route create

Create or update route

### Synopsis

Create or update route

```
knctl route create [flags]
```

### Examples

```

  # Set traffic percentages for service 'svc1' in namespace 'ns1'
  knctl route create --route rt1 -p svc1:latest=20% -p svc1:previous=80% -n ns1

  # Roll back traffic for previous revision of service 'svc1' in namespace 'ns1'
  knctl route create --route rt1 -p svc1:previous=100% -n ns1
```

### Options

```
      --generate-name        Set to generate name
  -h, --help                 help for create
  -n, --namespace string     Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -p, --percentage strings   Set percentage (format: revision=percentage, example: latest=100%) (can be specified multiple times)
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

* [knctl route](knctl_route.md)	 - Route management (create, curl, delete, list, show)

