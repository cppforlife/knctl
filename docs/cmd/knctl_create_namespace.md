## knctl create namespace

Create namespace

### Synopsis

Create namespace.

Use 'kubectl delete ns <name>' to delete namespace.

```
knctl create namespace [flags]
```

### Examples

```

  # Create namespace 'ns1'
  knctl create namespace -n ns1
```

### Options

```
  -h, --help               help for namespace
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

* [knctl create](knctl_create.md)	 - Create resources (namespace, service-account, basic-auth-secret, ssh-auth-secret)

