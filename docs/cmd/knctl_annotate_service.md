## knctl annotate service

Annotate service

### Synopsis

Annotate service

```
knctl annotate service [flags]
```

### Examples

```

  # Annotate service 'srv1' in namespace 'ns1' with key and value
  knctl annotate service -s srv1 -a key=value -n ns1
```

### Options

```
  -a, --annotation strings   Set annotation (format: key=value) (can be specified multiple times)
  -h, --help                 help for service
  -n, --namespace string     Specified namespace (or default from kubeconfig)
  -s, --service string       Specified service
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

* [knctl annotate](knctl_annotate.md)	 - Annotate (revision, service)

