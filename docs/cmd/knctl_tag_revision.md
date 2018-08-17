## knctl tag revision

Tag revision

### Synopsis

Tag revision

```
knctl tag revision [flags]
```

### Examples

```

  # Tag revision 'rev1' in namespace 'ns1' as 'stable'
  knctl tag revision -r rev1 -t stable -n ns1
```

### Options

```
  -h, --help               help for revision
  -n, --namespace string   Specified namespace (or default from kubeconfig)
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

* [knctl tag](knctl_tag.md)	 - Tag (revision)

