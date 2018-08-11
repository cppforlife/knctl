## knctl delete revision

Delete revision

### Synopsis

Delete revision

```
knctl delete revision [flags]
```

### Examples

```

  # Delete revision 'rev1' in namespace 'ns1'
  knctl delete revision -r rev1 -n ns1
```

### Options

```
  -h, --help               help for revision
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE)
  -r, --revision string    Specified revision
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

