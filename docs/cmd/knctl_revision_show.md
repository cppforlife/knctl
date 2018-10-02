## knctl revision show

Show revision

### Synopsis

Show revision details in a namespace

```
knctl revision show [flags]
```

### Examples

```

  # Show details for revison 'rev1' in namespace 'ns1'
  knctl revision show -r rev1 -n ns1
```

### Options

```
  -h, --help               help for show
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -r, --revision string    Specified revision
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

* [knctl revision](knctl_revision.md)	 - Revision management (annotate, delete, list, show, tag, untag)

