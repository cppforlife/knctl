## knctl namespace create

Create namespace

### Synopsis

Create namespace.

Use 'kubectl delete ns <name>' to delete namespace.

```
knctl namespace create [flags]
```

### Examples

```

  # Create namespace 'ns1'
  knctl namespace create -n ns1
```

### Options

```
      --generate-name      Set to generate name
  -h, --help               help for create
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
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

* [knctl namespace](knctl_namespace.md)	 - Namespace (create)

