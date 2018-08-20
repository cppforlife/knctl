## knctl create ssh-auth-secret

Create SSH auth secret

### Synopsis

Create SSH auth secret.

Use 'kubectl get secret -n <namespace>' to list secrets.
Use 'kubectl delete secret <name> -n <namespace>' to delete secret.

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
      --generate-name             Set to generate name
      --github                    Preconfigure type and url for Github.com Git access
  -h, --help                      help for ssh-auth-secret
      --known-hosts string        Set known hosts
  -n, --namespace string          Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
      --private-key string        Set private key in PEM format ($KNCTL_SSH_AUTH_SECRET_PRIVATE_KEY)
      --private-key-path string   Set private key in PEM format from file path
  -s, --secret string             Specified secret
      --type string               Set type (example: git)
      --url string                Set url (example: github.com)
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

