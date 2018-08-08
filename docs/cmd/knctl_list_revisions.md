## knctl list revisions

List revisions

### Synopsis

List all revisions for a service

```
knctl list revisions [flags]
```

### Examples

```

  # List all revisions for service 'svc1' in namespace 'ns1' 
  knctl list revisions -s svc1 -n ns1
```

### Options

```
  -h, --help               help for revisions
  -n, --namespace string   Specified namespace
  -s, --service string     Specified service
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

