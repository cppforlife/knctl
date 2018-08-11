## knctl delete build

Delete build

### Synopsis

Delete build

```
knctl delete build [flags]
```

### Examples

```

  # Delete build 'build1' in namespace 'ns1'
  knctl delete build -b build1 -n ns1
```

### Options

```
  -b, --build string       Specified build
  -h, --help               help for build
  -n, --namespace string   Specified namespace (can be provided via environment variable KNCTL_NAMESPACE)
```

### Options inherited from parent commands

```
      --column strings      Filter to show only given columns
      --json                Output as JSON
      --kubeconfig string   Path to the kubeconfig file (can be provided via environment variable KNCTL_KUBECONFIG) (default "/Users/pivotal/.kube/config")
      --no-color            Disable colorized output
      --non-interactive     Don't ask for user input
      --tty                 Force TTY-like output
```

### SEE ALSO

* [knctl delete](knctl_delete.md)	 - Delete resource (service, revision, route, build)

