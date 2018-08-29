## knctl service-account create

Create service account

### Synopsis

Create service account.

Use 'kubectl get serviceaccount -n <namespace>' to list service accounts.
Use 'kubectl delete serviceaccount <name> -n <namespace>' to delete service account.

```
knctl service-account create [flags]
```

### Examples

```

  # Create service account 'sa1' with two secrets in namespace 'ns1'
  knctl service-account create -a sa1 --secret secret1 --secret secret2 -n ns1
```

### Options

```
      --generate-name               Set to generate name
  -h, --help                        help for create
  -p, --image-pull-secret strings   Set image pull secret (format: secret-name) (can be specified multiple times)
  -n, --namespace string            Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -s, --secret strings              Set secret (format: secret-name) (can be specified multiple times)
  -a, --service-account string      Specified service-account
```

### Options inherited from parent commands

```
      --column strings              Filter to show only given columns
      --json                        Output as JSON
      --kubeconfig string           Path to the kubeconfig file ($KNCTL_KUBECONFIG or $KUBECONFIG)
      --kubeconfig-context string   Kubeconfig context override ($KNCTL_KUBECONFIG_CONTEXT)
      --no-color                    Disable colorized output
      --non-interactive             Don't ask for user input
      --tty                         Force TTY-like output
```

### SEE ALSO

* [knctl service-account](knctl_service-account.md)	 - Service account management (create)

