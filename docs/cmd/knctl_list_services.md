## knctl list services

List services

### Synopsis

List all services in a namespace

```
knctl list services [flags]
```

### Examples

```

  # List all services in namespace 'ns1'
  knctl list services -n ns1
```

### Options

```
  -h, --help               help for services
  -n, --namespace string   Specified namespace
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

* [knctl list](knctl_list.md)	 - List resources (service, revision, ingress)

