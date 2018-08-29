## knctl basic-auth-secret create

Create basic auth secret

### Synopsis

Create basic auth secret.

Use 'kubectl get secret -n <namespace>' to list secrets.
Use 'kubectl delete secret <name> -n <namespace>' to delete secret.

```
knctl basic-auth-secret create [flags]
```

### Examples

```

  # Create SSH basic auth secret 'secret1' in namespace 'ns1'
  knctl basic-auth-secret create -s secret1 --type ssh --url github.com --username username --password password -n ns1

  # Create Docker registry basic auth secret 'secret1' in namespace 'ns1'
  knctl basic-auth-secret create -s secret1 --docker-hub --username username --password password -n ns1

  # Create Docker registry basic auth secret 'secret1' for pulling images in namespace 'ns1'
  knctl basic-auth-secret create -s secret1 --docker-hub --username username --password password --for-pulling -n ns1

  # Create GCR.io registry basic auth secret 'secret1' in namespace 'ns1'
  knctl basic-auth-secret create -s secret1 --gcr --username username --password password -n ns1

  # Create generic Docker registry basic auth secret 'secret1' in namespace 'ns1'
  knctl basic-auth-secret create -s secret1 --type docker --url https://registry.domain.com/ --username username --password password -n ns1
```

### Options

```
      --docker-hub         Preconfigure type and url for Docker Hub registry
      --for-pulling        Convert to pull secret ('kubernetes.io/dockerconfigjson' type)
      --gcr                Preconfigure type and url for gcr.io registry
      --generate-name      Set to generate name
  -h, --help               help for create
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -p, --password string    Set password ($KNCTL_BASIC_AUTH_SECRET_PASSWORD)
  -s, --secret string      Specified secret
      --type string        Set type (example: docker, ssh)
      --url string         Set url (example: https://index.docker.io/v1/, https://github.com)
  -u, --username string    Set username
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

* [knctl basic-auth-secret](knctl_basic-auth-secret.md)	 - Basic auth secret management (create)

