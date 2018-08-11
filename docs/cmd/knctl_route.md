## knctl route

Configure route

### Synopsis

Configure route

```
knctl route [flags]
```

### Examples

```

  # Set traffic percentages for service 'svc1' in namespace 'ns1'
  knctl route --route rt1 -p svc1:latest=20% -p svc1:previous=80% -n ns1

  # Roll back traffic for previous revision of service 'svc1' in namespace 'ns1'
  knctl route --route rt1 -p svc1:previous=100% -n ns1
```

### Options

```
  -h, --help                 help for route
  -n, --namespace string     Specified namespace ($KNCTL_NAMESPACE)
  -p, --percentage strings   Set percentage (format: revision=percentage, example: latest=100%) (can be specified multiple times)
      --route string         Specified route
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

* [knctl](knctl.md)	 - knctl controls Knative resources (annotate, build, create, curl, delete, deploy, install, list, logs, open, route, tag, untag, version)

