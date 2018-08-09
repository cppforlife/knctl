## knctl create ssh-auth-secret

Create basic auth secret

### Synopsis

Create basic auth secret

```
knctl create ssh-auth-secret [flags]
```

### Examples

```

  # Create SSH secret 'secret1' in namespace 'ns1'
  knctl create ssh-auth-secret -s secret1 --url github.com --private-key ... --known-hosts ... -n ns1
```

### Options

```
  -h, --help                 help for ssh-auth-secret
      --known-hosts string   Set known hosts
  -n, --namespace string     Specified namespace (can be provided via environment variable KNCTL_NAMESPACE)
      --private-key string   Set private key in PEM format
  -s, --secret string        Specified secret
      --url string           Set url (example: https://github.com)
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

* [knctl create](knctl_create.md)	 - Create resources (namespace)

