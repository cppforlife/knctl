## knctl annotate revision

Annotate revision

### Synopsis

Annotate revision

```
knctl annotate revision [flags]
```

### Examples

```

  # Annotate revision 'rev1' in namespace 'ns1' with key and value
  knctl annotate revision -r rev1 -a key=value -n ns1
```

### Options

```
  -a, --annotation strings   Set annotation (format: key=value) (can be specified multiple times)
  -h, --help                 help for revision
  -n, --namespace string     Specified namespace
  -r, --revision string      Specified revision
```

### Options inherited from parent commands

```
      --column strings      Filter to show only given columns
      --json                Output as JSON
      --kubeconfig string   Path to the kubeconfig file (default "/Users/pivotal/.kube/config")
      --no-color            Disable colorized output
      --non-interactive     Don't ask for user input
      --tty                 Force TTY-like output
```

### SEE ALSO

* [knctl annotate](knctl_annotate.md)	 - Annotate resources (revision)

