## knctl create

Create resources (namespace, service-account, basic-auth-secret, ssh-auth-secret)

### Synopsis

Create resources (namespace, service-account, basic-auth-secret, ssh-auth-secret)

```
knctl create [flags]
```

### Options

```
  -h, --help   help for create
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

* [knctl](knctl.md)	 - knctl controls Knative resources
* [knctl create basic-auth-secret](knctl_create_basic-auth-secret.md)	 - Create basic auth secret
* [knctl create domain](knctl_create_domain.md)	 - Create domain
* [knctl create namespace](knctl_create_namespace.md)	 - Create namespace
* [knctl create service-account](knctl_create_service-account.md)	 - Create service account
* [knctl create ssh-auth-secret](knctl_create_ssh-auth-secret.md)	 - Create SSH auth secret

