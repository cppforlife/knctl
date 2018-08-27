## knctl revision annotate

Annotate revision

### Synopsis

Annotate revision

```
knctl revision annotate [flags]
```

### Examples

```

  # Annotate revision 'rev1' in namespace 'ns1' with key and value
  knctl revision annotate -r rev1 -a key=value -n ns1
```

### Options

```
  -a, --annotation strings   Set annotation (format: key=value) (can be specified multiple times)
  -h, --help                 help for annotate
  -n, --namespace string     Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -r, --revision string      Specified revision
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

* [knctl revision](knctl_revision.md)	 - Revision (annotate, delete, list, tag, untag)

