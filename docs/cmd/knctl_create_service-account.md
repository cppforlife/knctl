## knctl create service-account

Create service account

### Synopsis

Create service account.

Use 'kubectl get serviceaccount -n <namespace>' to list service accounts.
Use 'kubectl delete serviceaccount <name> -n <namespace>' to delete service account.

```
knctl create service-account [flags]
```

### Examples

```

  # Create service account 'sa1' with two secrets in namespace 'ns1'
  knctl create service-account -a sa1 --secret secret1 --secret secret2 -n ns1
```

### Options

```
      --generate-name               Set to generate name
  -h, --help                        help for service-account
  -p, --image-pull-secret strings   Set image pull secret (format: secret-name) (can be specified multiple times)
  -n, --namespace string            Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -s, --secret strings              Set secret (format: secret-name) (can be specified multiple times)
  -a, --service-account string      Specified service-account
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

* [knctl create](knctl_create.md)	 - Create (basic-auth-secret, domain, namespace, service-account, ssh-auth-secret)

