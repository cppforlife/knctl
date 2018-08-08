## knctl logs

Print logs

### Synopsis

Print logs of all active pods for a service

```
knctl logs [flags]
```

### Examples

```

  # Fetch last 10 log lines for service 'svc1' in namespace 'ns1' 
  knctl logs -s svc1 -n ns1

  # Follow logs for service 'svc1' in namespace 'ns1' 
  knctl logs -f -s svc1 -n ns1
```

### Options

```
  -f, --follow             As new revisions are added, new pod logs will be printed
  -h, --help               help for logs
  -l, --lines int          Number of lines (default 10)
  -n, --namespace string   Specified namespace (can be provided via environment variable KNCTL_NAMESPACE)
  -s, --service string     Specified service
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

* [knctl](knctl.md)	 - knctl controls Knative resources

