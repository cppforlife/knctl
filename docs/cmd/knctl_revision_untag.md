## knctl revision untag

Untag revision

### Synopsis

Untag revision

```
knctl revision untag [flags]
```

### Examples

```

  # Untag revision 'rev1' in namespace 'ns1' as 'stable'
  knctl revision untag -r rev1 -t stable -n ns1
```

### Options

```
  -h, --help               help for untag
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -r, --revision string    Specified revision
  -t, --tag strings        Set tag (format: value) (can be specified multiple times)
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

* [knctl revision](knctl_revision.md)	 - Revision (annotate, delete, list, tag, untag)

