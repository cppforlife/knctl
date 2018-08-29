## knctl revision tag

Tag revision

### Synopsis

Tag revision

```
knctl revision tag [flags]
```

### Examples

```

  # Tag revision 'rev1' in namespace 'ns1' as 'stable'
  knctl revision tag -r rev1 -t stable -n ns1
```

### Options

```
  -h, --help               help for tag
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -r, --revision string    Specified revision
  -t, --tag strings        Set tag (format: value) (can be specified multiple times)
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

* [knctl revision](knctl_revision.md)	 - Revision management (annotate, delete, list, tag, untag)

