## knctl delete service

Delete service

### Synopsis

Delete service

```
knctl delete service [flags]
```

### Examples

```

  # Delete service 'svc1' in namespace 'ns1'
  knctl delete service -s svc1 -n ns1
```

### Options

```
  -h, --help               help for service
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE)
  -s, --service string     Specified service
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

* [knctl delete](knctl_delete.md)	 - Delete resource (service, revision, route, build)

