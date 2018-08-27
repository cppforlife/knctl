## knctl revision list

List revisions

### Synopsis

List all revisions for a service

```
knctl revision list [flags]
```

### Examples

```

  # List all revisions for service 'svc1' in namespace 'ns1' 
  knctl revision list -s svc1 -n ns1
```

### Options

```
  -h, --help               help for list
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

* [knctl revision](knctl_revision.md)	 - Revision (annotate, delete, list, tag, untag)

