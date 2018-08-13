## knctl create basic-auth-secret

Create basic auth secret

### Synopsis

Create basic auth secret.

Use 'kubectl delete secret <name> -n <namespace>' to delete secret.

```
knctl create basic-auth-secret [flags]
```

### Examples

```

  # Create SSH basic auth secret 'secret1' in namespace 'ns1'
  knctl create basic-auth-secret -s secret1 --type ssh --url github.com --username username --password password -n ns1

  # Create Docker registry basic auth secret 'secret1' in namespace 'ns1'
  knctl create basic-auth-secret -s secret1 --docker-hub --username username --password password -n ns1

  # Create Docker registry basic auth secret 'secret1' for pulling images in namespace 'ns1'
  knctl create basic-auth-secret -s secret1 --docker-hub --username username --password password --for-pulling -n ns1

  # Create GCR.io registry basic auth secret 'secret1' in namespace 'ns1'
  knctl create basic-auth-secret -s secret1 --gcr --username username --password password -n ns1

  # Create generic Docker registry basic auth secret 'secret1' in namespace 'ns1'
  knctl create basic-auth-secret -s secret1 --type docker --url https://registry.domain.com/ --username username --password password -n ns1
```

### Options

```
      --docker-hub         Use Docker Hub registry (automatically fills 'type' and 'url')
      --for-pulling        Convert to pull secret ('kubernetes.io/dockerconfigjson' type)
      --gcr                Use gcr.io registry (automatically fills 'type' and 'url')
  -h, --help               help for basic-auth-secret
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE)
  -p, --password string    Set password ($KNCTL_BASIC_AUTH_SECRET_PASSWORD)
  -s, --secret string      Specified secret
      --type string        Set type (example: docker, ssh)
      --url string         Set url (example: https://index.docker.io/v1/, https://github.com)
  -u, --username string    Set username
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

