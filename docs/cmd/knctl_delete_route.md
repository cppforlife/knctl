## knctl delete route

Delete route

### Synopsis

Delete route

```
knctl delete route [flags]
```

### Examples

```

  # Delete route 'route1' in namespace 'ns1'
  knctl delete route --route route1 -n ns1
```

### Options

```
  -h, --help               help for route
  -n, --namespace string   Specified namespace (or default from kubeconfig)
      --route string       Specified route
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

* [knctl delete](knctl_delete.md)	 - Delete (build, revision, route, service)

