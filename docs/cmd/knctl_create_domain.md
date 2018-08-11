## knctl create domain

Create domain

### Synopsis

Create domain

```
knctl create domain [flags]
```

### Examples

```

  # Create domain 'example.com' and set it as default
  knctl create domain -d example.com --default
```

### Options

```
      --default         Set domain as default (currently required to be provided)
  -d, --domain string   Specified domain (example: domain.com)
  -h, --help            help for domain
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

