## knctl untag revision

Untag revision

### Synopsis

Untag revision

```
knctl untag revision [flags]
```

### Examples

```

  # Untag revision 'rev1' in namespace 'ns1' as 'stable'
  knctl untag revision -r rev1 -t stable -n ns1
```

### Options

```
  -h, --help               help for revision
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE)
  -r, --revision string    Specified revision
  -t, --tag strings        Set tag (format: value) (can be specified multiple times)
```

### Options inherited from parent commands

```
      --column strings      Filter to show only given columns
      --json                Output as JSON
      --kubeconfig string   Path to the kubeconfig file ($KNCTL_KUBECONFIG) (default "/Users/pivotal/.kube/config")
      --no-color            Disable colorized output
      --non-interactive     Don't ask for user input
      --tty                 Force TTY-like output
```

### SEE ALSO

* [knctl untag](knctl_untag.md)	 - Untag (revision)

