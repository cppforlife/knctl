## knctl route curl

Curl route

### Synopsis

Send a HTTP request to the first ingress address with the Host header set to the service's domain.

Requires 'curl' command installed on the system.

```
knctl route curl [flags]
```

### Examples

```

  # Curl route 'rt1' in namespace 'ns1'
  knctl route curl --route rt1 -n ns1
```

### Options

```
  -h, --help               help for curl
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -p, --port int32         Set port (default 80)
      --route string       Specified route
  -v, --verbose            Makes curl verbose during the operation
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

